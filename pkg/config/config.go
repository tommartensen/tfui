package config

import "github.com/tommartensen/tfui/pkg/util"

type Configuration struct {
	ServerToken string
	BaseDir     string
	Port        string
	Addr        string
	ClientToken string
}

// Initialize the Configuration struct by reading the values for it from the environment variables.
func New() *Configuration {
	return &Configuration{
		ServerToken: util.GetEnv("APPLICATION_TOKEN", ""),
		BaseDir:     util.GetEnv("BASE_DIR", "./plans"),
		Port:        util.GetEnv("PORT", "8080"),
		Addr:        util.GetEnv("TFUI_ADDR", "http://localhost:8080"),
		ClientToken: util.GetEnv("TFUI_TOKEN", ""),
	}
}
