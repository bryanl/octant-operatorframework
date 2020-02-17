package main

import (
	"log"

	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"

	"github.com/bryanl/octant-operatorframework/pkg/oof"
)

var pluginName = "operator-framework"

func main() {
	capabilities := &plugin.Capabilities{
		IsModule: true,
	}

	options := []service.PluginOption{
		service.WithNavigation(oof.HandleNavigation, oof.InitRoutes),
	}

	p, err := service.Register(pluginName, "Operator Framework", capabilities, options...)
	if err != nil {
		log.Fatal(err)
	}

	p.Serve()
}
