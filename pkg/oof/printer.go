package oof

import (
	"fmt"

	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

var (
	emptyResponse = plugin.PrintResponse{}
)

type ObjectPrinter interface {
	PrintObject(request *service.PrintRequest) (plugin.PrintResponse, error)
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
	default:
		return emptyResponse, fmt.Errorf("unable to print object of type %s", groupVersionKind)
	}

	return objectPrinter.PrintObject(request)
}
