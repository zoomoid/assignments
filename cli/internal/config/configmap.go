package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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
	config := &Configuration{}
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	in, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	err = Unmarshal(in, config)
	return config, err
}

// Write marshals a configuration struct into YAML and writes it to the designated file
func Write(config *Configuration, path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	out, err := Marshal(*config)
	if err != nil {
		return err
	}
	_, err = fd.Write(out)
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
