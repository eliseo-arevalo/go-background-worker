package internal

import (
    "encoding/json"
    "os"
    "log"
)

type APIConfig struct {
    URL      string `json:"url"`
    Interval string `json:"interval"`
}

func LoadConfig() ([]APIConfig, error) {
    apisEnv := os.Getenv("APIS")
    if apisEnv == "" {
        log.Println("⚠️ APIS environment variable not found, using default configuration")
        return []APIConfig{}, nil
    }

    var apiConfigs []APIConfig
    if err := json.Unmarshal([]byte(apisEnv), &apiConfigs); err != nil {
        return nil, err
    }

    return apiConfigs, nil
}