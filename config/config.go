package config

var (
	Config map[string]bool
	// dry-run recursive
	DataDir   string // Config file Directory
	OutputDir string // Output file Directory
	Proxy     string // Github proxy url
)

func init() {
	Config = make(map[string]bool)
	DataDir = "data"
	OutputDir = "."
	Proxy = ""
}
