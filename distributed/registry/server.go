/*
	Package registry provides the regisgry service.

	It's used in a distributed system and works
	as a registry server. It helps to locate
	specific services in the network.

	Workflow:

	The registry server maintains a database recording
	information about online services in the network.
	It accepts registry request and sends corresponding
	notifications. Specifically,

	If service A is online, it will send register request
	to the registry server and enters the registration
	database.

	Other online services (that required this service)
	will be notified of this service info. Besides,
	this service will be notified of info of its
	"required" services.

	If service A is offline, it will send deregister request
	to the registry server, and corresponding notifications
	will be sent out.

	The registry server also maintains a goroutine (heartbeat)
	periodically checking status of services in its database.
	It will remove any offline service from the database, and
	then notify corresponding services (that require it).

	Structure:

	--------------                          ------------------
	|            |  <--- (de)register ----  |                | ------
	|            |  ----    notify    --->  |    service 1   |      |
	|            |  ----   heartbeat  --->  | (require 2, x) | ---  |
	|            |                          ------------------   |  |
	|            |                                               |  |
	|            |                          ------------------   |  |
	|            |  <--- (de)register ----  |                | <--  |
	|            |  ----    notify    --->  |    service 2   |      |
	|            |  ----   heartbeat  --->  |   (require y)  |      |
	|            |                          ------------------      |
	|  Registry  |                                   .              |
	|            |                                   .              |
	|   Service  |                                   .              |
	|            |                          ------------------      |
	|            |  <--- (de)register ----  |                | <-----
	|            |  ----    notify    --->  |    service x   |
	|            |  ----   heartbeat  --->  |    (require N) | ----
	|            |                          ------------------    |
	|            |                                                |
	|            |                          ------------------    |
	|            |  <--- (de)register ----  |                |    |
	|            |  ----              ----  |    service N   | <---
	|            |  ----   heartbeat  --->  |                |
	--------------                          ------------------

	Protocols: http 1.1

	How to register a service?

	This registry package also provides apis to (de)register
	a service. The client service just need to call these
	apis for a (de)registration:

		func RegisterService(s Registration) error
		func DeregisterService(s Registration) error

*/
package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	ServiceName = "Registry Service"
	ServiceURL  = "http://localhost:3000"
)

const (
	Endpoint_Register   = "/register"
	Endpoint_Deregister = "/deregister"
	Endpoint_Services   = "/services"
	Endpoint_HeartBeat  = "/heartbeat"
)

func RegisterHandlers() {
	http.HandleFunc(Endpoint_Register, makeHandler(Endpoint_Register))
	http.HandleFunc(Endpoint_Deregister, makeHandler(Endpoint_Deregister))
}

const (
	HeartBeatCount         = 3
	HeartBeatCountInterval = 10 * time.Second
	HeartBeatInterval      = 60 * time.Second
)

// Services is a database maintaining online services
type Services struct {
	rs map[string]ServiceInfo
	sync.Mutex
}

var services = &Services{
	rs: make(map[string]ServiceInfo),
}

func (p *Services) add(s ServiceInfo) {
	p.Lock()
	defer p.Unlock()
	p.rs[s.Name] = s
	log.Printf("add registration:%v", s)
	log.Printf("current services:%v", p.rs)
	p.sendRequiredService(s)
	p.notify(patch{
		Added: []patchEntry{{s.Name, s.URL}},
	})
}

func (p *Services) remove(s ServiceInfo) {
	p.Lock()
	defer p.Unlock()
	delete(p.rs, s.Name)
	log.Printf("remove registration:%v", s)
	log.Printf("current services:%v", p.rs)
	p.notify(patch{
		Removed: []patchEntry{{s.Name, s.URL}},
	})
}

func (srvs *Services) sendRequiredService(s ServiceInfo) {
	var p patch
	var needUpdate = false
	for k, v := range srvs.rs {
		for _, required := range s.Required {
			if k == required {
				p.Added = append(p.Added, patchEntry{v.Name, v.URL})
				needUpdate = true
			}
		}
	}
	if needUpdate {
		go p.sendPatch(s.UpdateURL)
	}
}

func (srvs *Services) notify(pt patch) {
	for _, s := range srvs.rs {
		go func(s ServiceInfo) {
			p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
			needUpdate := false
			for _, required := range s.Required {
				for _, entry := range pt.Added {
					if required == entry.Name {
						p.Added = append(p.Added, entry)
						needUpdate = true
					}
				}
				for _, entry := range pt.Removed {
					if required == entry.Name {
						p.Removed = append(p.Removed, entry)
						needUpdate = true
					}
				}
			}
			if needUpdate {
				go p.sendPatch(s.UpdateURL)
			}
		}(s)
	}
}

func (srvs *Services) heartBeatOnce() {
	var wg sync.WaitGroup
	for _, s := range srvs.rs {
		wg.Add(1)
		go func(s ServiceInfo) {
			lastIsSuccess := true
			defer wg.Done()
			for i := 0; i < HeartBeatCount; i++ {
				_, err := do(struct{}{}, s.HeartBeat, "")
				// failed -- remove service
				if err != nil {
					log.Printf("heartbeat test %v fail", s.Name)
					if lastIsSuccess {
						log.Printf("heartbeat remove %v", s.Name)
						srvs.remove(s)
						lastIsSuccess = false
					}
					time.Sleep(HeartBeatCountInterval)
					continue
				}
				// pass
				log.Printf("heartbeat test %v pass", s.Name)
				if !lastIsSuccess {
					log.Printf("heartbeat add %v", s.Name)
					srvs.add(s)
				}
				break
			}
		}(s)
	}
	wg.Wait()
}

func StartHeartBeat() {
	var once sync.Once
	once.Do(func() {
		go func() {
			for {
				services.heartBeatOnce()
				time.Sleep(HeartBeatInterval)
			}
		}()
	})
}

func makeHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var s ServiceInfo
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch path {
		case Endpoint_Register:
			services.add(s)
		case Endpoint_Deregister:
			services.remove(s)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
