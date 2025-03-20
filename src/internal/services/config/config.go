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

// Service is the interface for the config service.
type Service interface {
	// GetConfig returns the config.
	GetConfig() *Config

	// LoadConfig loads the config from the config file.
	LoadConfig() error

	// SaveConfig saves the config to the config file.
	SaveConfig(*Config) error

	// SetConfigFromUI sets the config from the UI.
	SetConfigFromUI(*Config)
}

type configService struct {
	config     *Config
	configDate time.Time
	configMu   sync.RWMutex
}

// NewService creates a new config service.
func NewService() Service {
	return &configService{}
}

func (c *configService) GetConfig() *Config {
	if c.configDate.IsZero() || time.Since(c.configDate).Minutes() > 1 {
		_ = c.LoadConfig()
		c.configDate = time.Now()
	}

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

func (c *configService) SaveConfig(cfg *Config) error {
	c.configMu.Lock()
	defer c.configMu.Unlock()

	cfgFile, err := os.OpenFile(filepath.Join(utils.IO.GetMainFolder(), "Config", "config.json"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	enc := json.NewEncoder(cfgFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return err
	}

	c.config = cfg
	c.configDate = time.Now()

	return nil
}

func (c *configService) SetConfigFromUI(cfg *Config) {
	c.configMu.Lock()
	defer c.configMu.Unlock()

	c.config = cfg
	c.configDate = time.Now()
}

func (c *configService) setupDefaultConfig() error {
	data, err := assets.ReadFile(assets.AssetDefaultConfig)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c.config)
}
