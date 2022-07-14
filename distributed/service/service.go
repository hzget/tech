/*
	Package service provide an interface to set and start a http server.

	Its users are web services that can make use of these common code
	to start a http server. This makes it easy to manage and develop
	these services.

	[*] register/deregister this service to the registry server automatically
	[*] the service developers focus on the service code -- http handlers

	The Start() function starts the server, create 2 goroutines and returns a context.
	The caller can block with <-ctx.Done() until the context channel is closed.

	If the user enters a newline in the standard input, the server will exit.
	And the context channel will be closed.

	An example:

	func main() {
		// arguments are provided by the user (the specific service)
		ctx, err := service.Start(ctx, *Registraion{xx,xx}, RegisterHandlersFunc)
		if err != nil { os.Exit(1) }
		<- ctx.Done()
	}

*/
package service

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Start() start the http server
// Arguments:
//	registry.Registraion - contains service info to register
//  RegisterHandlersFunc - http handlers - the service code
func Start(ctx context.Context, s registry.Registration,
	RegisterHandlersFunc func()) (context.Context, error) {
	RegisterHandlersFunc()
	return startService(ctx, s)
}

func startService(ctx context.Context, s registry.Registration) (context.Context, error) {

	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	url, err := url.Parse(s.URL)
	if err != nil {
		return ctx, err
	}

	srv.Addr = url.Host

	log.Printf("[%v] at [%v] started, enter a newline to stop it", s.Name, srv.Addr)
	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		var s1 string
		fmt.Scanln(&s1)
		registry.DeregisterService(s)
		srv.Shutdown(ctx)
		cancel()
	}()

	registry.RegisterService(s)

	return ctx, nil
}
