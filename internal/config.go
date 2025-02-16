package internal

import (
    "encoding/json"
    "os"

    "github.com/joho/godotenv"
)

type APIConfig struct {
    URL      string `json:"url"`
    Interval string `json:"interval"`
}

func LoadConfig() ([]APIConfig, error) {
    if err := godotenv.Load(); err != nil {
        return nil, err
    }

    apisEnv := os.Getenv("APIS")
    if apisEnv == "" {
        return nil, nil
    }

    var apiConfigs []APIConfig
    if err := json.Unmarshal([]byte(apisEnv), &apiConfigs); err != nil {
        return nil, err
    }

    return apiConfigs, nil
}