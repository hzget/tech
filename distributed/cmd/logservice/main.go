package main

import (
	"context"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	stdlog "log"
)

func main() {
	log.Run("./logfile.txt")
	s := registry.ServiceInfo{
		Name:      log.ServiceName,
		URL:       log.ServiceURL,
		UpdateURL: log.ServiceURL + "/services",
		HeartBeat: log.ServiceURL + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(), s, log.RegisterHandlers)
	if err != nil {
		stdlog.Println(err)
	}

	<-ctx.Done()

	stdlog.Printf("service [%v] exit", log.ServiceName)
}
