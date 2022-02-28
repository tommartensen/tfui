package config

import "github.com/tommartensen/tfui/pkg/util"

type Configuration struct {
	ApplicationToken string
	BaseDir          string
	Port             string
}

// Initialize the Configuration struct by reading the values for it from the environment variables.
func New() *Configuration {
	return &Configuration{
		ApplicationToken: util.GetEnv("APPLICATION_TOKEN", ""),
		BaseDir:          util.GetEnv("BASE_DIR", "./plans"),
		Port:             util.GetEnv("PORT", "8080"),
	}
}
