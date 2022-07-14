/*
	Package log provides "log" service:
		accept "http post" request and write its body to a log file

	endpoint: /log
*/
package log

import (
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
)

const (
	ServiceName = "Log Service"
	ServiceURL  = "http://localhost:4000"
)

const (
	Endpoint_Logging = "/log"
)

const (
	LogPrefix = "[go]: "
)

// It will write message to an io.Writer object - a filename
// which implements an io.Writer interface
var log *stdlog.Logger

// Run create a new Logger whose destination is a file
// It shall be called firstly.
func Run(filename string) {
	log = stdlog.New(filelog(filename), LogPrefix, stdlog.LstdFlags|stdlog.Lmicroseconds)
}

type filelog string

func (fl filelog) Write(b []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Write(b)
}

// RegisterHandlers:
//
// for the endpoint "/log", register a handler
// to accept "http post" request and write its body to a log file
//
// It's the keypoint of the log service.
func RegisterHandlers() {
	http.HandleFunc(Endpoint_Logging, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil || len(body) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// keypoint
			log.Printf(string(body))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
