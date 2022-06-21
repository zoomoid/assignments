package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

// ReadConfigMap reads in a config file an unmarshals it into a configuration struct
func ReadConfigMap(path string) (*Configuration, error) {
	if path == "" {
		path = "."
	}
	viper.SetConfigFile(".assignments.yaml")
	viper.SetConfigType("yaml")
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

// WriteConfigMap converts a configuration struct into a map for viper to write to the filesystem as config file
func WriteConfigMap(config *Configuration, path string) error {
	if path == "" {
		path = "."
	}
	viper.SetConfigFile(".assignments.yaml")
	viper.SetConfigType("yaml")
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
