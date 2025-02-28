package config

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/assets"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/go-json-experiment/json/v1"
)

// ConfigService is the interface for the config service.
type ConfigService interface {
	// GetConfig returns the config.
	GetConfig() *Config

	// LoadConfig loads the config from the config file.
	LoadConfig() error
}

// Implementation check
var _ ConfigService = (*configService)(nil)

type configService struct {
	config   *Config
	configMu sync.RWMutex
}

// NewConfigService creates a new config service.
func NewConfigService() *configService {
	return &configService{}
}

func (c *configService) GetConfig() *Config {
	c.configMu.RLock()
	defer c.configMu.RUnlock()

	return c.config
}

func (c *configService) LoadConfig() error {
	c.configMu.Lock()
	defer c.configMu.Unlock()

	cfgFile, err := os.Open(filepath.Join(utils.IO.GetMainFolder(), "Config", "config.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return c.setupDefaultConfig()
		}

		return err
	}
	defer cfgFile.Close()

	if err := json.NewDecoder(cfgFile).Decode(c.config); err != nil {
		return err
	}

	time.Local = c.config.GetTimezone()

	return nil
}

func (c *configService) setupDefaultConfig() error {
	data, err := assets.ReadFile(assets.AssetDefaultConfig)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c.config)
}
