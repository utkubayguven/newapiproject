package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Config struct {
	APIPort           int       `json:"api_port"`
	MaxRequestsPerDay int       `json:"max_requests_per_day"`
	RemainingRequests int       `json:"remaining_requests"`
	LastReset         time.Time `json:"last_reset"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		file, err := os.Open("config/config.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		instance = &Config{}
		err = decoder.Decode(instance)
		if err != nil {
			panic(err)
		}

		resetIfNeeded(instance)
	})
	return instance
}

func resetIfNeeded(c *Config) {
	now := time.Now()
	if now.Sub(c.LastReset) >= 24*time.Hour {
		c.RemainingRequests = c.MaxRequestsPerDay
		c.LastReset = now
		saveConfig(c)
	}
}

func saveConfig(c *Config) {
	file, err := os.Create("config/config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		panic(err)
	}
}

func (c *Config) DecreaseRequestCount() error {
	resetIfNeeded(c)

	if c.RemainingRequests > 0 {
		c.RemainingRequests--
		saveConfig(c)
		return nil
	}
	return fmt.Errorf("daily request limit exceeded")
}
