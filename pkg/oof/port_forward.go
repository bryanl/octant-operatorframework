package oof

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//go:generate mockgen -destination=./fake/mock_port_forward.go -package=fake github.com/bryanl/octant-operatorframework/pkg/oof PortForward

type PortForward interface {
	ToService(ctx context.Context, namespace, name string, port uint16) (string, func(), error)
}

// DashboardPortForward manages port forwards.
type DashboardPortForward struct {
	Dashboard service.Dashboard
}

// NewDashboardPortForward creates an instance of DashboardPortForward.
func NewDashboardPortForward(dashboard service.Dashboard) (PortForward, error) {
	if dashboard == nil {
		return nil, fmt.Errorf("dashboard client is nil")
	}

	pf := DashboardPortForward{
		Dashboard: dashboard,
	}

	return &pf, nil
}

// ToService creates a port forward to a service.
func (pf *DashboardPortForward) ToService(ctx context.Context, namespace, name string, port uint16) (string, func(), error) {
	endpoints, err := pf.getEndpoints(ctx, namespace, name)
	if err != nil {
		return "", nil, fmt.Errorf("get endpoints: %w", err)
	}

	pod, err := pf.endpointPod(ctx, endpoints)
	if err != nil {
		return "", nil, fmt.Errorf("find pod in endpoints: %w", err)
	}

	req := api.PortForwardRequest{
		Namespace: namespace,
		PodName:   pod.GetName(),
		Port:      port,
	}
	resp, err := pf.Dashboard.PortForward(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("create port forward")
	}

	cancel := func() {
		pf.Dashboard.CancelPortForward(ctx, resp.ID)
	}

	uri := fmt.Sprintf("localhost:%d", resp.Port)
	logrus.WithField("address", uri).Info("created port forward")

	return uri, cancel, nil
}

// endpointPod retrieves a pod from endpoints. It chooses the target of the first address of the first
// subset. If unable, it returns an error. This will also return an error if the target ref
// in the address does not point to a pod.
func (pf *DashboardPortForward) endpointPod(ctx context.Context, endpoints *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	rawSubsets, _, err := unstructured.NestedSlice(endpoints.Object, "subsets")
	if err != nil {
		return nil, fmt.Errorf("get subsets from from endpoints: %w", err)
	}

	if len(rawSubsets) < 1 {
		return nil, fmt.Errorf("there were no subsets defined")
	}

	subset, ok := rawSubsets[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unknown subset format")
	}

	rawAddresses, _, err := unstructured.NestedSlice(subset, "addresses")
	if err != nil {
		return nil, fmt.Errorf("get addresses from subset: %w", err)
	}

	if len(rawAddresses) < 1 {
		return nil, fmt.Errorf("there were no addresses defined")
	}

	address, ok := rawAddresses[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unknown subset format")
	}

	kind, _, err := unstructured.NestedString(address, "targetRef", "kind")
	if err != nil {
		return nil, fmt.Errorf("get target kind: %w", err)
	}

	if kind != "Pod" {
		return nil, fmt.Errorf("unable to handle %s", kind)
	}

	name, _, err := unstructured.NestedString(address, "targetRef", "name")
	if err != nil {
		return nil, fmt.Errorf("get target name: %w", err)
	}

	return pf.getPod(ctx, endpoints.GetNamespace(), name)
}

func (pf *DashboardPortForward) getEndpoints(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	endpointsKey := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Endpoints",
		Name:       name,
	}

	return pf.Dashboard.Get(ctx, endpointsKey)
}

func (pf *DashboardPortForward) getPod(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	podKey := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       name,
	}

	return pf.Dashboard.Get(ctx, podKey)
}
