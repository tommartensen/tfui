package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tommartensen/tfui/pkg/config"
)

func TestNewDefaultConfiguration(t *testing.T) {
	c := config.New()
	assert.Equal(t, c.ApplicationToken, "")
	assert.Equal(t, c.BaseDir, "./plans")
	assert.Equal(t, c.Port, "8080")
}

func TestConfigurationWithEnvironment(t *testing.T) {
	os.Setenv("APPLICATION_TOKEN", "testtoken")
	defer os.Unsetenv("APPLICATION_TOKEN")
	os.Setenv("BASE_DIR", "./test/plans")
	defer os.Unsetenv("BASE_DIR")
	os.Setenv("PORT", "4000")
	defer os.Unsetenv("PORT")

	c := config.New()
	assert.Equal(t, c.ApplicationToken, "testtoken")
	assert.Equal(t, c.BaseDir, "./test/plans")
	assert.Equal(t, c.Port, "4000")
}
