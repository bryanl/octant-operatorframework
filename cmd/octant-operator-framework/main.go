package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	sigCh := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		done <- true
	}()

	go func() {
		logrus.Info("serve plugin")
		p.Serve()

	}()

	<-done
	logrus.Info("plugin is exiting")
}
