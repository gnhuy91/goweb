package vcapparser

import "encoding/json"

type Credentials struct {
	ID         int    `json:"ID"`
	BindingID  string `json:"binding_id"`
	Database   string `json:"database"`
	DSN        string `json:"dsn"`
	Host       string `json:"host"`
	InstanceID string `json:"instance_id"`
	JdbcURI    string `json:"jdbc_uri"`
	Password   string `json:"password"`
	Port       string `json:"port"`
	URI        string `json:"uri"`
	Username   string `json:"username"`
}

type VcapService struct {
	Credentials    Credentials `json:"credentials"`
	Label          string      `json:"label"`
	Name           string      `json:"name"`
	Plan           string      `json:"plan"`
	Provider       string      `json:"provider"`
	SyslogDrainURL string      `json:"syslog_drain_url"`
	Tags           []string    `json:"tags"`
}

// VcapServices is a map of services detail
type VcapServices map[string][]VcapService

// ParseVcapServices parse string provided from VCAP_SERVICES environment var
// to VcapServices struct.
func ParseVcapServices(vcapStr string) (VcapServices, error) {
	var vcapServices VcapServices

	if err := json.Unmarshal([]byte(vcapStr), &vcapServices); err != nil {
		return vcapServices, err
	}

	return vcapServices, nil
}
