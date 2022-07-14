package main

import (
	"context"
	"distributed/registry"
	"distributed/service"
	"log"
)

func main() {
	registry.StartHeartBeat()
	s := registry.ServiceInfo{
		Name:      registry.ServiceName,
		URL:       registry.ServiceURL,
		UpdateURL: registry.ServiceURL + "/services",
	}
	ctx, err := service.Start(context.Background(), s, registry.RegisterHandlers)
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	log.Printf("service [%v] exit", registry.ServiceName)
}
