package oof

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"k8s.io/apimachinery/pkg/runtime"
)

//go:generate mockgen -destination=./fake/mock_print_request.go -package=fake github.com/bryanl/octant-operatorframework/pkg/oof PrintRequest
//go:generate mockgen -destination=./fake/mock_dashboard_client.go -package=fake github.com/vmware-tanzu/octant/pkg/plugin/service Dashboard

var (
	emptyResponse = plugin.PrintResponse{}
)

type PrintRequest interface {
	Context() context.Context
	Object() runtime.Object
	DashboardClient() service.Dashboard
}

type printRequest struct {
	request *service.PrintRequest
}

func newPrintRequest(request *service.PrintRequest) *printRequest {
	pr := printRequest{request: request}
	return &pr
}

var _ PrintRequest = &printRequest{}

func (p printRequest) Context() context.Context {
	return p.request.Context()
}

func (p printRequest) Object() runtime.Object {
	return p.request.Object
}

func (p printRequest) DashboardClient() service.Dashboard {
	return p.request.DashboardClient
}

type ObjectPrinter interface {
	PrintObject(request PrintRequest) (plugin.PrintResponse, error)
}

type Printer struct{}

func NewPrinter() *Printer {
	p := Printer{}

	return &p
}

func (p *Printer) HandlePrint(request *service.PrintRequest) (plugin.PrintResponse, error) {

	if request.Object == nil {
		return plugin.PrintResponse{}, fmt.Errorf("unable to print a nil object")
	}

	groupVersionKind := request.Object.GetObjectKind().GroupVersionKind()

	var objectPrinter ObjectPrinter

	switch groupVersionKind {
	case CatalogSourceGVK:
		objectPrinter = NewCatalogSourcePrinter()
	case SubscriptionGVK:
		objectPrinter = NewSubscriptionPrinter()
	default:
		return emptyResponse, fmt.Errorf("unable to print object of type %s", groupVersionKind)
	}

	pr := newPrintRequest(request)

	resp, err := objectPrinter.PrintObject(pr)
	if err != nil {
		logrus.WithError(err).Error("unable to print object")
		return emptyResponse, fmt.Errorf("print object")
	}

	return resp, nil
}
