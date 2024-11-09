package sbom

import (
	"bytes"
	"os"
	"strings"

	cdx "github.com/CycloneDX/cyclonedx-go"
)

func NewCycloneDx(file string) (*cdx.BOM, error) {
	// parse sbom
	bom := cdx.NewBOM()

	// read data
	data, err := os.ReadFile(file)
	if err != nil {
		return bom, err
	}

	// infer sbom format
	var format cdx.BOMFileFormat
	if strings.HasSuffix(file, ".xml") {
		format = cdx.BOMFileFormatXML
	} else if strings.HasSuffix(file, ".json") {
		format = cdx.BOMFileFormatJSON
	}

	decoder := cdx.NewBOMDecoder(bytes.NewReader(data), format)
	return bom, decoder.Decode(bom)
}
