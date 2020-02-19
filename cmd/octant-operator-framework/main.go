package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"

	"github.com/bryanl/octant-operatorframework/pkg/oof"
)

var pluginName = "operator-framework"

func main() {
	filename := "/tmp/of.log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	logrus.SetOutput(f)

	defer func() {
		if r := recover(); r != nil {
			logrus.Println("recovered", r)
		}
	}()

	printer := oof.NewPrinter()

	options := []service.PluginOption{
		service.WithPrinter(printer.HandlePrint),
		service.WithNavigation(oof.HandleNavigation, oof.InitRoutes),
	}

	logrus.Info("registering service")
	p, err := service.Register(pluginName, "Operator Framework", oof.InitCapabilities(), options...)
	if err != nil {
		log.Fatal(err)
	}

	logrus.Info("serving")
	p.Serve()
}
