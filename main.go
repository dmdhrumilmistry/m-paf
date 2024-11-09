package main

import (
	"flag"

	_ "github.com/dmdhrumilmistry/m-pat/pkg/logging"
	"github.com/dmdhrumilmistry/m-pat/pkg/sbom"
	"github.com/rs/zerolog/log"
)

func main() {
	filePath := flag.String("f", "", "path of sbom file")
	flag.Parse()

	bom, err := sbom.NewCycloneDx(*filePath)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to parse SBOM file: %s", *filePath)
	}

	for _, component := range *bom.Components {
		log.Print(component.PackageURL)
	}
}
