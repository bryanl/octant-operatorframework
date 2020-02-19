package oof

import (
	"fmt"
	"strconv"

	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type CatalogSourcePrinterOption func(CatalogSourcePrinter) CatalogSourcePrinter

type CatalogSourcePrinter struct {
	RegistryClientFactory func(address string) (RegistryClient, error)
	PortForwardFactory    func(dashboard service.Dashboard) (PortForward, error)
}

var _ ObjectPrinter = (*CatalogSourcePrinter)(nil)

func NewCatalogSourcePrinter(options ...CatalogSourcePrinterOption) *CatalogSourcePrinter {
	csp := CatalogSourcePrinter{}

	for _, option := range options {
		csp = option(csp)
	}

	if csp.RegistryClientFactory == nil {
		csp.RegistryClientFactory = NewGRPCRegistryClient
	}

	if csp.PortForwardFactory == nil {
		csp.PortForwardFactory = NewDashboardPortForward
	}

	return &csp
}

func (c CatalogSourcePrinter) PrintObject(request PrintRequest) (plugin.PrintResponse, error) {
	if request == nil {
		return emptyResponse, fmt.Errorf("unable to print a nil object")
	}

	object, err := toUnstructured(request.Object())
	if err != nil {
		return emptyResponse, fmt.Errorf("invalid object: %w", err)
	}

	response := plugin.PrintResponse{}

	if response.Config, err = addStringSectionItem(response.Config, "Image", object.Object,
		"spec", "image"); err != nil {
		return emptyResponse, err
	}

	if response.Config, err = addStringSectionItem(response.Config, "Source Type", object.Object,
		"spec", "sourceType"); err != nil {
		return emptyResponse, err
	}

	if response.Status, err = addStringSectionItem(response.Status, "Connection State",
		object.Object, "status", "connectionState", "lastObservedState"); err != nil {
		return emptyResponse, err
	}

	packageTable, err := c.printPackages(request, object.Object)
	if err != nil {
		return emptyResponse, err
	}

	response.Items = append(response.Items, component.FlexLayoutItem{
		Width: component.WidthFull,
		View:  packageTable,
	})

	return response, nil
}

var (
	CatalogSourcePackageCols = component.NewTableCols("Name")
)

func (c *CatalogSourcePrinter) printPackages(request PrintRequest, m map[string]interface{}) (*component.Table, error) {
	sp, err := c.catalogServicePort(m)
	if err != nil {
		return nil, fmt.Errorf("get catalog registry address: %w", err)
	}

	pf, err := c.PortForwardFactory(request.DashboardClient())
	if err != nil {
		return nil, fmt.Errorf("create port forward: %w", err)
	}
	uri, cancel, err := pf.ToService(request.Context(), sp.Namespace, sp.Name, sp.Port)
	if err != nil {
		return nil, fmt.Errorf("create port forward to service: %w", err)
	}
	defer cancel()

	registryClient, err := c.RegistryClientFactory(uri)
	if err != nil {
		return nil, fmt.Errorf("create registry client: %w", err)
	}

	packageTable := component.NewTable("Packages", "", CatalogSourcePackageCols)

	packageNames, err := registryClient.ListPackages(request.Context())
	if err != nil {
		return nil, err
	}

	for _, name := range packageNames {
		row := component.TableRow{
			"Name": component.NewText(name),
		}
		packageTable.Add(row)
	}

	return packageTable, nil
}

type unknownRegistryError struct {
	protocol string
}

var _ error = &unknownRegistryError{}

func newUnknownRegistryError(protocol string) *unknownRegistryError {
	return &unknownRegistryError{protocol: protocol}
}

func (e *unknownRegistryError) Error() string {
	return fmt.Sprintf("unknown catalog registry protocol: %q", e.protocol)
}

type servicePort struct {
	Namespace string
	Name      string
	Port      uint16
}

func (c *CatalogSourcePrinter) catalogServicePort(m map[string]interface{}) (servicePort, error) {
	protocol, _, err := unstructured.NestedString(m, "status", "registryService", "protocol")
	if err != nil {
		return servicePort{}, err
	}

	if protocol != "grpc" {
		return servicePort{}, newUnknownRegistryError(protocol)
	}

	portStr, _, err := unstructured.NestedString(m, "status", "registryService", "port")
	if err != nil {
		return servicePort{}, err
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return servicePort{}, err
	}

	name, _, err := unstructured.NestedString(m, "status", "registryService", "serviceName")
	if err != nil {
		return servicePort{}, err
	}

	namespace, _, err := unstructured.NestedString(m, "status", "registryService", "serviceNamespace")
	if err != nil {
		return servicePort{}, err
	}

	return servicePort{
		Namespace: namespace,
		Name:      name,
		Port:      uint16(port),
	}, nil
}

func toUnstructured(in runtime.Object) (*unstructured.Unstructured, error) {
	if in == nil {
		return nil, fmt.Errorf("object is nil")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(in)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: m}, nil
}

func addStringSectionItem(sections []component.SummarySection, name string, m map[string]interface{}, fields ...string) ([]component.SummarySection, error) {
	text, _, err := unstructured.NestedString(m, fields...)
	if err != nil {
		return nil, err
	}

	return append(sections, component.SummarySection{
		Header:  name,
		Content: component.NewText(text),
	}), nil
}
