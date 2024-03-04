package config

import "os"

type Config struct {
	path       string
	boolFields map[string]field[bool]
}

func NewConfig(configPath string) *Config {
	c := Config{
		path:       configPath,
		boolFields: make(map[string]field[bool]),
	}
	return &c
}

func (c *Config) AddBool(fieldName string, required bool, defaultValue bool) {
	c.boolFields[fieldName] = field[bool]{Default: defaultValue, Required: required}
}

// Read() reads the config file.
func (c *Config) Read() error {
	return nil
}

// Flush() flushes all the changes.
func (c *Config) Flush() error {
	if err := os.MkdirAll(c.path, os.ModePerm); err != nil {
		return err
	}

	return nil
}

type field[T any] struct {
	Default  T
	Required bool
}
