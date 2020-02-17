package oof

import (
	"fmt"

	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type CatalogSourcePrinter struct{}

var _ ObjectPrinter = (*CatalogSourcePrinter)(nil)

func NewCatalogSourcePrinter() *CatalogSourcePrinter {
	csp := CatalogSourcePrinter{}

	return &csp
}

func (c CatalogSourcePrinter) PrintObject(request *service.PrintRequest) (plugin.PrintResponse, error) {
	if request == nil {
		return emptyResponse, fmt.Errorf("unable to print a nil object")
	}

	object, err := toUnstructured(request.Object)
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

	return response, nil
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
