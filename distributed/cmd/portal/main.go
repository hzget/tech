package main

import (
	"context"
	"distributed/log"
	"distributed/portal"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
	"time"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	host, port := "localhost", "5000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		Name: registry.PortalService,
		URL:  serviceAddress,
		Required: []string{
			registry.LogService,
			registry.GradingService,
		},
		UpdateURL: serviceAddress + "/services",
		HeartBeat: serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(),
		r,
		portal.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	time.Sleep(time.Second)
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s\n", logProvider)
		log.SetClientLogger(logProvider, r.Name)
	} else {
		fmt.Printf("[%v] not found, err: %v\n", registry.LogService, err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down portal.")
}
