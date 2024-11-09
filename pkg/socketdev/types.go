package socketdev

// Start of Alert Types

type AlertDescription struct {
	Emoji         string            `json:"emoji"`
	Props         map[string]string `json:"props"`
	Title         string            `json:"title"`
	Suggestion    string            `json:"suggestion"`
	Description   string            `json:"description"`
	NextStepTitle string            `json:"nextStepTitle"`
}

type AlertType struct {
	Id        int                         `json:"id"`
	CreatedAt string                      `json:"created_at"`
	UpdatedAt string                      `json:"updated_at"`
	Type      string                      `json:"type"`
	I18n      map[string]AlertDescription `json:"i18n"`
	Category  string                      `json:"category"`
	Severity  float32                     `json:"severity"`
	Enabled   bool                        `json:"enabled"`
}

type ResAlertTypes map[string]AlertType
type AlertTypes map[int]AlertType

// End of Alert Types

// Start of Package Info Types
type Alert struct {
	Type int    `json:"type"`
	Key  string `json:"key"`
}
type Qualifiers struct {
	Ext string `json:"ext"`
}

type Scores struct {
	SupplyChain   float64 `json:"supplyChain"`
	Quality       float64 `json:"quality"`
	Maintenance   float64 `json:"maintenance"`
	Vulnerability float64 `json:"vulnerability"`
	License       float64 `json:"license"`
	Overall       float64 `json:"overall"`
}

type Capabilities struct {
	Env    bool `json:"env"`
	Eval   bool `json:"eval"`
	Fs     bool `json:"fs"`
	Net    bool `json:"net"`
	Shell  bool `json:"shell"`
	Unsafe bool `json:"unsafe"`
}
type PackageInfo struct {
	ID             string       `json:"id"`
	Type           string       `json:"type"`
	Name           string       `json:"name"`
	Namespace      string       `json:"namespace"`
	Files          string       `json:"files"`
	Version        string       `json:"version"`
	Qualifiers     Qualifiers   `json:"qualifiers"`
	Scores         Scores       `json:"scores"`
	Capabilities   Capabilities `json:"capabilities"`
	License        string       `json:"license"`
	Size           int          `json:"size"`
	State          string       `json:"state"`
	Alerts         []Alert      `json:"alerts"`
	LicenseDetails []any        `json:"licenseDetails"`
}

// End of Package Info Types
