package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config application configuration
type Config struct {
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		SSL      bool   `json:"sslmode"`
	} `json:"database_settings"`
	Logger struct {
		DumpRequest string `json:"dump_request"`
	} `json:"logger_settings"`
}

// InitConfig parse configuration file and setup settings
func InitConfig(filePath string) (*Config, error) {
	var config Config
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &config, err
	}
	json.Unmarshal(file, &config)
	return &config, nil
}
