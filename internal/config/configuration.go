package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration contains the details of the application configuration.
type Configuration struct {
	// filePath is the path to the configuration file. This file contains the
	// application data in JSON format.
	filePath string
	// Data is the application data. This data is loaded from the configuration
	// file and is used to configure the application.
	Data *ApplicationData
}

// NewConfiguration returns a new Configuration from a JSON file.
func NewConfiguration(filePath string) (*Configuration, error) {
	configuration := &Configuration{
		filePath: filePath,
		Data:     newDefaultApplicationData(),
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			_ = configuration.Flush() // ignore error if file write failed
			return configuration, nil
		}
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(configuration.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return configuration, nil
}

// Flush serializes and writes the configuration to the file.
func (c *Configuration) Flush() error {
	file, err := os.Create(c.filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c.Data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return file.Sync()
}
