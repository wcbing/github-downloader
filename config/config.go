package config

var (
	Config map[string]bool
	// dry-run dir proxy recursive
	DataDir string // Config file Directory
	Proxy   string // Github proxy url
)

func init() {
	Config = make(map[string]bool)
	DataDir = "data"
}
