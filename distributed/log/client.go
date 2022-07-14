package log

import (
	"bytes"
	"fmt"
	stdlog "log"
	"net/http"
)

func SetClientLogger(url string, prefix string) {
	stdlog.SetPrefix("[" + prefix + "]: ")
	stdlog.SetFlags(0)
	stdlog.SetOutput(clientLogger(url + Endpoint_Logging))
}

type clientLogger string

// Write(p []byte) (n int, err error)
func (cl clientLogger) Write(p []byte) (n int, err error) {

	client := &http.Client{}
	resp, err := client.Post(string(cl), "plain/text", bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%v return status [%d] is not %d",
			string(cl), resp.StatusCode, http.StatusOK)
	}

	return len(p), nil
}
