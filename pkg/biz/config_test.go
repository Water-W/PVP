package biz

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	c.Worker.Addresses = []string{
		"1.2.3.4:1203",
		"2.3.4.5:1204",
	}
	c.DumpTo(os.Stdout)
}

func TestLoad(t *testing.T) {
	c, err := LoadConfigFromFile("../../master.config.example.toml")
	assert.NoError(t, err)
	c.DumpTo(os.Stdout)
}
