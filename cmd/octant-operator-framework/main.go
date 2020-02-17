package main

import (
	"log"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"

	"github.com/bryanl/octant-operatorframework/pkg/oof"
)

var pluginName = "operator-framework"

func main() {
	printer := oof.NewPrinter()

	options := []service.PluginOption{
		service.WithPrinter(printer.HandlePrint),
		service.WithNavigation(oof.HandleNavigation, oof.InitRoutes),
	}

	p, err := service.Register(pluginName, "Operator Framework", oof.InitCapabilities(), options...)
	if err != nil {
		log.Fatal(err)
	}

	p.Serve()
}
