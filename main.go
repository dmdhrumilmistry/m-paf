package main

import (
	"flag"
	"fmt"
	"slices"
	"strings"

	_ "github.com/dmdhrumilmistry/m-paf/pkg/logging"
	"github.com/dmdhrumilmistry/m-paf/pkg/sbom"
	"github.com/dmdhrumilmistry/m-paf/pkg/socketdev"
	"github.com/rs/zerolog/log"
)

func banner() {
	fmt.Print(`
======================================================

      ███╗   ███╗      ██████╗  █████╗ ███████╗
      ████╗ ████║      ██╔══██╗██╔══██╗██╔════╝
      ██╔████╔██║█████╗██████╔╝███████║█████╗  
      ██║╚██╔╝██║╚════╝██╔═══╝ ██╔══██║██╔══╝  
      ██║ ╚═╝ ██║      ██║     ██║  ██║██║     
      ╚═╝     ╚═╝      ╚═╝     ╚═╝  ╚═╝╚═╝     
------------------------------------------------------
            Malicious-PAckageFinder
======================================================
  Detect malicious packages and risks from SBOM file
------------------------------------------------------
         github.com/dmdhrumilmistry/m-paf
------------------------------------------------------

`)
}

func main() {
	banner()

	filePath := flag.String("f", "", "Path to CycloneDX SBOM file")
	workers := flag.Int("w", 50, "Number of workers for scanning malicious package in SBOM. Default value: 50")
	acceptableSCRisk := flag.Float64("t", 0.5, "Acceptable Supply Chain Risk Score Threshold. Higher the score better, lower the risk. Ranges between 0 to 1. Default: 0.5")
	flag.Parse()

	api, err := socketdev.NewSocketAPI()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to fetch alert types")
	}

	bom, err := sbom.NewCycloneDx(*filePath)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to parse SBOM file: %s", *filePath)
	}

	packagesInfo := api.ProcessComponents(bom, *workers)

	scRiskyPackages := []string{}
	alertsMap := map[int][]string{}
	for _, packageInfo := range packagesInfo {
		// detect packages with lower supply chain security score
		if packageInfo.Scores.SupplyChain < *acceptableSCRisk {
			scRiskyPackages = append(scRiskyPackages, fmt.Sprintf("%s (%.2f)", packageInfo.Name, packageInfo.Scores.SupplyChain))
		}

		// create map of alerts
		for _, alert := range packageInfo.Alerts {
			if !slices.Contains(alertsMap[alert.Type], packageInfo.Name) {
				alertsMap[alert.Type] = append(alertsMap[alert.Type], packageInfo.Name)
			}
		}
	}

	log.Info().Msg("Package Analysis Completed")
	log.Info().Msg("Packages having supply chain risk lower than acceptable value")
	fmt.Println("Risky Packages: " + strings.Join(scRiskyPackages, ", "))
	fmt.Println("======================================================")

	log.Info().Msg("Alerts for Packages")
	for alertId, packages := range alertsMap {
		alert := api.AlertTypes[alertId].I18n["en-US"]
		fmt.Printf("%s - %s\n", alert.Emoji, alert.Title)
		fmt.Println("Description:", alert.Description)
		fmt.Println("Suggestion:", alert.Suggestion)
		fmt.Println("Risky Packages:", strings.Join(packages, ", "))
		fmt.Println("------------------------------------------------------")
	}
	fmt.Println("======================================================")
}
