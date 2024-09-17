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
	Orphans_Ignore   Orphans = "ignore"
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

	config := &Config{
		Upload:  true,
		Orphans: Orphans_Ignore,
		Include: make([]string, 0, 1),
		Exclude: make([]string, 0, 1),
		Run:     make([]string, 0, 1),
	}

	err = json.Unmarshal(data, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
