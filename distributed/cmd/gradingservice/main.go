package main

import (
	"context"
	"distributed/grades"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stdlog "log"
	"time"
)

func main() {

	r := registry.ServiceInfo{
		Name:      grades.ServiceName,
		URL:       grades.ServiceURL,
		Required:  []string{registry.LogService},
		UpdateURL: grades.ServiceURL + "/services",
		HeartBeat: grades.ServiceURL + "/heartbeat",
	}
	ctx, err := service.Start(context.Background(), r, grades.RegisterHandlers)
	if err != nil {
		stdlog.Fatal(err)
	}

	time.Sleep(time.Second)
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s\n", logProvider)
		log.SetClientLogger(logProvider, r.Name)
	} else {
		fmt.Printf("[%v] not found, err: %v\n", registry.LogService, err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
