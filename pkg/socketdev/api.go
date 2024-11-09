package socketdev

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	packageurl "github.com/package-url/packageurl-go"
)

type Api struct {
	BaseUrl    string
	AlertTypes AlertTypes

	// Endpoints
	PollWithAlertsEndpoint string
	GetAlertTypesEndpoint  string
}

func NewSocketAPI() (*Api, error) {
	api := &Api{
		BaseUrl:                "https://socket.dev/api/",
		PollWithAlertsEndpoint: "ecosystems/artifact/poll-with-alerts",
		GetAlertTypesEndpoint:  "ecosystems/alert/alert-types",
	}

	if err := api.getAlertTypes(); err != nil {
		return nil, err
	}

	return api, nil
}

/*
Fetches risk alerts for provided package url (purl)
example: api.GetAlerts("pkg:maven/org.dependencytrack/dependency-track@4.1.0")
*/
func (c *Api) GetAlerts(purl string) ([]PackageInfo, error) {
	var packageInfo []PackageInfo

	// Build the base URL
	endpoint, err := url.Parse(c.BaseUrl + c.PollWithAlertsEndpoint)
	if err != nil {
		return nil, err
	}

	// parse package Url
	pUrl, err := packageurl.FromString(purl)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	queryParams := url.Values{}
	queryParams.Add("name", pUrl.Name)
	queryParams.Add("namespace", pUrl.Namespace)
	queryParams.Add("type", pUrl.Type)
	queryParams.Add("version", pUrl.Version)

	endpoint.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned %d instead of 200, with resp: %v", res.StatusCode, res)
	}

	// Use bufio.Scanner to read the response line-by-line
	scanner := bufio.NewScanner(res.Body)
	if scanner.Scan() {
		var packageDetail PackageInfo
		line := scanner.Text()

		// Decode the first line as JSON into PackageDetails
		if err := json.Unmarshal([]byte(line), &packageDetail); err != nil {
			return nil, fmt.Errorf("failed to decode JSON: %w", err)
		}

		packageInfo = append(packageInfo, packageDetail)
	}

	return packageInfo, nil
}

func (c *Api) getAlertTypes() error {
	var alertTypes ResAlertTypes
	url := c.BaseUrl + c.GetAlertTypesEndpoint

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if !slices.Contains([]int{200, 304}, res.StatusCode) {
		return fmt.Errorf("api returned %d, expected 200/304 - %v", res.StatusCode, res)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &alertTypes); err != nil {
		return err
	}

	alertMaps := make(AlertTypes)
	for _, alertType := range alertTypes {
		alertMaps[alertType.Id] = alertType
	}

	c.AlertTypes = alertMaps

	return nil
}
