package main

import (
	"flag"

	_ "github.com/dmdhrumilmistry/m-pat/pkg/logging"
	"github.com/dmdhrumilmistry/m-pat/pkg/sbom"
	"github.com/dmdhrumilmistry/m-pat/pkg/socketdev"
	"github.com/rs/zerolog/log"
)

func main() {
	filePath := flag.String("f", "", "path of sbom file")
	flag.Parse()

	api, err := socketdev.NewSocketAPI()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to fetch alert types")
	}

	bom, err := sbom.NewCycloneDx(*filePath)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to parse SBOM file: %s", *filePath)
	}

	for _, component := range *bom.Components {

		if component.PackageURL != "pkg:maven/org.apache.poi/poi-ooxml-schemas@3.17?type=jar" {
			continue
		}

		log.Info().Msgf("Processing purl component - %s", component.PackageURL)
		packageInfo, err := api.GetAlerts(component.PackageURL)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get alerts for purl: %s", component.PackageURL)
			continue
		}
		log.Info().Interface("package info", packageInfo).Msg("")

		// for _, alertData := range packageInfo[0].Alerts {
		// 	alert := api.AlertTypes[alertData.Type].I18n["en-US"]
		// 	msg := fmt.Sprintf("%s - %s\nDescription: %s\nSuggestion: %s", alert.Emoji, alert.Title, alert.Description, alert.Suggestion)
		// 	log.Warn().Msg(msg)
		// }
		// log.Print("----------------------------")
	}

}
