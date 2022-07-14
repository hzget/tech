package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

type patchEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type patch struct {
	Added   []patchEntry `json:"added"`
	Removed []patchEntry `json:"removed"`
}

func (p patch) sendPatch(url string) error {
	_, err := do(p, url, "")
	if err != nil {
		fmt.Printf("fail to send patch: %v, error: %v\n", p, err)
	}
	return err
}

type providers struct {
	service map[string][]string
	sync.RWMutex
}

var pvds = &providers{
	service: make(map[string][]string),
}

func (p *providers) update(pt patch) {
	p.Lock()
	defer p.Unlock()
	for _, entry := range pt.Added {
		if _, ok := p.service[entry.Name]; !ok {
			p.service[entry.Name] = make([]string, 0)
		}
		p.service[entry.Name] = append(p.service[entry.Name], entry.URL)
	}
	for _, entry := range pt.Removed {
		addrs, ok := p.service[entry.Name]
		if !ok {
			continue
		}
		for i := range addrs {
			if addrs[i] == entry.URL {
				// i+1 overflows the slice?
				p.service[entry.Name] = append(addrs[:i], addrs[i+1:]...)
				break
			}
		}
		if len(p.service[entry.Name]) == 0 {
			delete(p.service, entry.Name)
		}
	}
}

func (p *providers) get(name string) (string, error) {
	p.RLock()
	defer p.RUnlock()
	addrs, ok := p.service[name]
	if !ok {
		return "", fmt.Errorf("service %s not found", name)
	}
	n := rand.Intn(len(addrs))
	return addrs[n], nil
}

func GetProvider(name string) (string, error) {
	return pvds.get(name)
}

type serviceUpdatehandler struct{}

func (uh serviceUpdatehandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var p patch
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("receive patch: %v\n", p)
	pvds.update(p)
	fmt.Printf("current providers: %v\n", pvds.service)
}

func addHeartBeatHandler(s Registration) error {
	heartbeatURL, err := url.Parse(s.HeartBeat)
	if err != nil {
		return err
	}
	http.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return nil
}

func addServiceUpdateHandler(s Registration) error {
	if s.Required != nil {
		updateurl, err := url.Parse(s.UpdateURL)
		if err != nil {
			return err
		}
		http.Handle(updateurl.Path, serviceUpdatehandler{})
	}
	return nil
}

func RegisterService(s Registration) error {
	if s.Name == ServiceName {
		return nil
	}
	// register http handlers for updating providers
	if err := addServiceUpdateHandler(s); err != nil {
		return err
	}
	// register http handlers for heartbeat
	if err := addHeartBeatHandler(s); err != nil {
		return err
	}
	// register this service
	_, err := do(s, ServiceURL, Endpoint_Register)
	return err
}

func DeregisterService(s Registration) error {
	if s.Name == ServiceName {
		return nil
	}
	_, err := do(s, ServiceURL, Endpoint_Deregister)
	return err
}

/*
func FindService(name string) (string, error) {
	resp, err := do(s, Endpoint_Find)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), err
}
*/

//func do[V Registration | patch](s V, url, req string) (*http.Response, error) {
func do(s interface{}, url, req string) (*http.Response, error) {
	buffer := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buffer).Encode(s); err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Post(url+req, "application/json", buffer)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("failed: [%v] return status %d", url, resp.StatusCode)
	}
	return resp, nil
}
