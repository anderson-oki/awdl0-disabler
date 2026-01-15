package configuration

import (
	"encoding/json"
	"os"
	"time"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
)

type JSONConfigAdapter struct {
	FilePath string
}

func NewJSONConfigAdapter(path string) *JSONConfigAdapter {
	return &JSONConfigAdapter{FilePath: path}
}

func (a *JSONConfigAdapter) Load() (*domain.Config, error) {
	file, err := os.Open(a.FilePath)
	if os.IsNotExist(err) {
		// Default Configuration
		return &domain.Config{PollingInterval: 1 * time.Second}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config domain.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	config.Clamp()

	return &config, nil
}

func (a *JSONConfigAdapter) Save(config *domain.Config) error {
	file, err := os.Create(a.FilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(config)
}
