package main

import (
	"fmt"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	"github.com/pivotal-cf/brokerapi"

	. "bosh-generic-sb/broker"
	btask "bosh-generic-sb/broker/task"
)

func main() {
	cfg, err := NewConfigFromPath(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	uuidGen := boshuuid.NewGenerator()
	logger := boshlog.NewLogger(boshlog.LevelInfo)

	directorImpl := NewDirectorImpl(cfg.Director)
	taskService := btask.NewAsyncTaskService(uuidGen, logger)
	apiImpl := NewBrokerAPIImpl(NewBrokerImpl(cfg.Broker, directorImpl), taskService)

	creds := brokerapi.BrokerCredentials{Username: cfg.HTTP.Username, Password: cfg.HTTP.Password}

	apiLogger := lager.NewLogger("bosh-generic-sb")
	apiLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	apiLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	http.Handle("/", brokerapi.New(apiImpl, apiLogger, creds))

	err = http.ListenAndServe(cfg.HTTP.Host+":"+cfg.HTTP.Port, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
