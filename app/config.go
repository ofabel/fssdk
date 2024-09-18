package app

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ofabel/fssdk/contract"
)

func getConfigFromFile(path string) (*contract.Config, error) {
	fp, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fp.Close()

	data, err := io.ReadAll(fp)

	if err != nil {
		return nil, err
	}

	config := &contract.Config{
		Upload:  true,
		Orphans: contract.Orphans_Ignore,
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
