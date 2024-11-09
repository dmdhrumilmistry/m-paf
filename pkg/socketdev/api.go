package socketdev

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"sync"

	cdx "github.com/CycloneDX/cyclonedx-go"
	packageurl "github.com/package-url/packageurl-go"
	"github.com/rs/zerolog/log"
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

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("api returned 404 for purl - %s", pUrl)
	}

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

func (c *Api) ProcessComponents(bom *cdx.BOM, workers int) []PackageInfo {

	jobQueue := make(chan cdx.Component, len(*bom.Components))
	results := make(chan PackageInfo, len(*bom.Components))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for component := range jobQueue {
				log.Info().Msgf("Processing purl component - %s", component.PackageURL)
				packageInfo, err := c.GetAlerts(component.PackageURL)
				if err != nil {
					log.Printf("failed to get alerts for purl: %s, error: %v", component.PackageURL, err)
					continue
				}

				results <- packageInfo[0]
			}
		}()
	}

	// Feed jobs to jobQueue
	for _, component := range *bom.Components {
		jobQueue <- component
	}
	close(jobQueue) // No more jobs

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	var packagesInfo []PackageInfo
	for result := range results {
		packagesInfo = append(packagesInfo, result)
	}

	return packagesInfo
	// Print all collected results at the end
	// for result := range results {
	// 	log.Printf("Alerts for %s", result.Name)
	// 	for _, alertData := range result.Alerts {
	// 		alert := c.AlertTypes[alertData.Type].I18n["en-US"]
	// 		msg := fmt.Sprintf("%s - %s\nDescription: %s\nSuggestion: %s", alert.Emoji, alert.Title, alert.Description, alert.Suggestion)
	// 		log.Print(msg)
	// 	}
	// 	log.Print("----------------------------")
	// }

	// for _, component := range *bom.Components {

	// 	if component.PackageURL != "pkg:maven/org.apache.poi/poi-ooxml-schemas@3.17?type=jar" {
	// 		continue
	// 	}

	// 	log.Info().Msgf("Processing purl component - %s", component.PackageURL)
	// 	packageInfo, err := c.GetAlerts(component.PackageURL)
	// 	if err != nil {
	// 		log.Error().Err(err).Msgf("failed to get alerts for purl: %s", component.PackageURL)
	// 		continue
	// 	}
	// 	log.Info().Interface("package info", packageInfo).Msg("")

	// 	for _, alertData := range packageInfo[0].Alerts {
	// 		alert := c.AlertTypes[alertData.Type].I18n["en-US"]
	// 		msg := fmt.Sprintf("%s - %s\nDescription: %s\nSuggestion: %s", alert.Emoji, alert.Title, alert.Description, alert.Suggestion)
	// 		log.Warn().Msg(msg)
	// 	}
	// 	log.Print("----------------------------")
	// }
}
