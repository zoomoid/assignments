package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ConfigurationFileName string = ".assignments.yaml"
	ConfigurationFileType string = "yaml"
)

var (
	ErrNoConfigmap error = errors.New("failed to find configmap in working directory or above")
)

// Read reads in a config file an unmarshals it into a configuration struct
func Read(path string) (*Configuration, error) {
	if path == "" {
		path = "."
	}
	viper.SetConfigFile(ConfigurationFileName)
	viper.SetConfigType(ConfigurationFileType)
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file, %v", err)
	}

	config := &Configuration{}

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config with viper, %v", err)
	}

	return config, nil
}

// Write converts a configuration struct into a map for viper to write to the filesystem as config file
func Write(config *Configuration, path string) error {
	if path == "" {
		path = "."
	}
	viper.SetConfigFile(ConfigurationFileName)
	viper.SetConfigType(ConfigurationFileType)
	viper.AddConfigPath(path)

	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var serializedConfig map[string]interface{}
	err = json.Unmarshal(marshalledConfig, &serializedConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to map, %v", err)
	}
	viper.MergeConfigMap(serializedConfig)

	isset := viper.IsSet("spec.course")
	if !isset {
		return fmt.Errorf("failed to marshal configuration for viper")
	}
	err = viper.WriteConfig()
	return err
}

// Find traverses the file system tree from the start path upwards in search of a
// configmap file.
//
// Returns the directory where we found the file. If the root is reached without
// finding a configmap, an error is returned
func Find(start string) (string, error) {
	p := start
	for {
		_, err := os.Stat(filepath.Join(p, ConfigurationFileName))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if p == filepath.Join(filepath.Dir(p), "..") {
					// path does not change by traversal, already at root
					break
				}
			}
			// move up one level and try again
			p = filepath.Join(p, "..")
			continue
		}
		// found configmap file at p
		return p, nil
	}
	return "", ErrNoConfigmap
}

// Unmarshal implements yaml unmarshalling for configuration structs
func Unmarshal(in []byte, out *Configuration) error {
	return yaml.Unmarshal(in, out)
}

// Marshal implements yaml marshalling for configuration structs
func Marshal(in Configuration) (out []byte, err error) {
	return yaml.Marshal(in)
}
