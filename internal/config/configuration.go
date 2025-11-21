package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/reugn/gemini-cli/gemini"
)

// Configuration contains the details of the application configuration.
type Configuration struct {
	// filePath is the path to the configuration file. This file contains the
	// application data in JSON format.
	filePath string
	// Data is the application data. This data is loaded from the configuration
	// file and is used to configure the application.
	Data *ApplicationData
	// lastModTime is the time the configuration file was last modified.
	lastModTime time.Time
}

// NewConfiguration returns a new Configuration from a JSON file.
// If the file does not exist, it is created with default values.
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
	if err := decoder.Decode(configuration.Data); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return configuration, nil
}

// Flush serializes and writes the configuration to the file.
//
// If the file is modified since the last load/flush, the configuration is re-read
// and merged with the on-disk data.
func (c *Configuration) Flush() error {
	// Reload the configuration if the file was modified since the last load/flush.
	if err := c.reloadIfStale(); err != nil {
		return err
	}

	// Create the file if it does not exist.
	file, err := os.Create(c.filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Serialize the configuration data to the file.
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c.Data); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	// Sync the file to disk.
	if err := file.Sync(); err != nil {
		return fmt.Errorf("error syncing file: %w", err)
	}

	// Get the file information.
	info, err := os.Stat(c.filePath)
	if err != nil {
		return fmt.Errorf("error stating file: %w", err)
	}

	// Update the last modified time.
	c.lastModTime = info.ModTime()

	return nil
}

// reloadIfStale re-reads and merges the on-disk configuration if the file was
// modified since the last load/flush.
func (c *Configuration) reloadIfStale() error {
	info, err := os.Stat(c.filePath)
	if err != nil {
		if os.IsNotExist(err) { // ignore error if file does not exist
			return nil
		}
		return fmt.Errorf("error stating config file: %w", err)
	}

	// If the file was not modified since the last load/flush, do nothing.
	if !c.lastModTime.IsZero() && !info.ModTime().After(c.lastModTime) {
		return nil
	}

	// Re-read the configuration file.
	file, err := os.Open(c.filePath)
	if err != nil {
		return fmt.Errorf("error reopening config file: %w", err)
	}
	defer file.Close()

	onDisk := newDefaultApplicationData()
	if err := json.NewDecoder(file).Decode(onDisk); err != nil {
		return fmt.Errorf("error decoding updated config: %w", err)
	}

	// Merge the on-disk data into the current configuration.
	c.mergeApplicationData(onDisk)

	return nil
}

// mergeApplicationData merges on-disk data into the current configuration.
func (c *Configuration) mergeApplicationData(onDisk *ApplicationData) {
	if onDisk == nil || c.Data == nil {
		return
	}

	// The CLI never modifies these fields; always overwrite with the on-disk values.
	c.Data.SystemPrompts = onDisk.SystemPrompts
	c.Data.SafetySettings = onDisk.SafetySettings
	c.Data.Tools = onDisk.Tools

	// Merge history records.
	if onDisk.History != nil {
		if c.Data.History == nil {
			c.Data.History = make(map[string][]*gemini.SerializableContent, len(onDisk.History))
		}
		for label, records := range onDisk.History {
			if _, exists := c.Data.History[label]; !exists {
				c.Data.History[label] = records
			}
		}
	}
}
