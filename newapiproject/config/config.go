package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	APIPort         int `json:"api_port"`
	MaxRequestLimit int `json:"max_request_limit"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		file, err := os.Open("config.json")
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
	})
	return instance
}
