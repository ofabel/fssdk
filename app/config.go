package app

import (
	"encoding/json"
	"io"
	"os"
)

type Orphans string

const (
	Orphans_Download Orphans = "download"
	Orphans_Delete   Orphans = "delete"
)

type Config struct {
	Source  string
	Target  string
	Upload  bool
	Orphans Orphans
	Include []string
	Exclude []string
	Run     []string
}

func GetConfigFromFile(path string) (*Config, error) {
	fp, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fp.Close()

	data, err := io.ReadAll(fp)

	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = json.Unmarshal(data, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
