package registry

// service info to register
type ServiceInfo struct {
	Name      string   `json:"name"`
	URL       string   `json:"url"`       // url to provide its service
	Required  []string `json:"required"`  // Required Services
	UpdateURL string   `json:"updateurl"` // url to accept update info
	HeartBeat string   `json:"heartbeat"` // url to accept heartbeat test
}

type Registration = ServiceInfo

const (
	LogService     = "Log Service"
	GradingService = "Grade Service"
	PortalService  = "Portal Service"
)
