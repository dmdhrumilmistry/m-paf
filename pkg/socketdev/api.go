package socketdev

import (
	"net/http"
	"net/url"
)

type Api struct {
	BaseUrl string

	// Endpoints
	PollWithAlertsEndpoint string
	GetAlertTypesEndpoint  string
}

func NewSocketAPI() *Api {
	return &Api{
		BaseUrl:                "https://socket.dev/api/",
		PollWithAlertsEndpoint: "ecosystems/artifact/poll-with-alerts",
		GetAlertTypesEndpoint:  "ecosystems/alert/alert-types",
	}
}

func (c *Api) GetAlerts(name, namespace, ptype, version string) (*http.Response, error) {
	// Build the base URL
	endpoint, err := url.Parse(c.BaseUrl + c.PollWithAlertsEndpoint)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	queryParams := url.Values{}
	if name != "" {
		queryParams.Add("name", name)
	}
	if namespace != "" {
		queryParams.Add("namespace", namespace)
	}
	if ptype != "" {
		queryParams.Add("ptype", ptype)
	}
	if version != "" {
		queryParams.Add("version", version)
	}
	endpoint.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

// func (c *Api)
