package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestEnvOutSide(t *testing.T) {
	t.Skip()
	viper.BindEnv("env")
	// c := viper.GetString("env")
	c := os.Getenv("env")
	assert.Equal(t, "dev", c)
}
