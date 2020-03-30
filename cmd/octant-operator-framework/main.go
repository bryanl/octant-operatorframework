package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"

	"github.com/bryanl/octant-operatorframework/pkg/oof"
)

var pluginName = "operator-framework"

func withLogger(fn func(logger logrus.FieldLogger)) {
	filename := "/tmp/of.log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	logger := logrus.New()
	logger.SetOutput(f)

	entry := logger.WithField("instance", uuid.New().String())

	fn(entry)
}

func main() {
	withLogger(func(logger logrus.FieldLogger) {
		logrus.Info("plugin is starting")

		defer func() {
			if r := recover(); r != nil {
				logrus.Println("recovered", r)
			}
		}()

		printer := oof.NewPrinter(logger)

		options := []service.PluginOption{
			service.WithPrinter(printer.Print),
			service.WithNavigation(oof.HandleNavigation, oof.InitRoutes),
			service.WithActionHandler(oof.HandleActions),
		}

		logger.Info("registering service")
		p, err := service.Register(pluginName, "Operator Framework", oof.InitCapabilities(), options...)
		if err != nil {
			log.Fatal(err)
		}

		sigCh := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigCh
			logger.WithField("signal", sigToString(sig)).Info("detected signal; preparing to exit")
			done <- true
		}()

		go func() {
			logger.Info("plugin is serving requests")
			p.Serve()

		}()

		<-done
		logger.Info("plugin is exiting")
	})
}

func sigToString(sig os.Signal) string {
	switch sig {
	case syscall.SIGINT:
		return "INT"
	case syscall.SIGTERM:
		return "TERM"
	default:
		return "unknown"
	}
}
