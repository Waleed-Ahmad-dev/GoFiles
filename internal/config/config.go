package config

import (
	"encoding/json"
	"os"
	"time"

	"GoFiles/internal/types"
)

// Global Config
const RootFolder = "."
const TrashFolder = ".trash"
const ThumbsFolder = ".thumbs" // NEW: Hidden folder for thumbnails
const TrashRetention = 30 * 24 * time.Hour
const ConfigFileName = "gofiles.json"

// Runtime State
var AppConfig types.ConfigFile
var IsConfigured = false

// InitConfig tries to load gofiles.json
func InitConfig() {
	file, err := os.Open(ConfigFileName)
	if err != nil {
		IsConfigured = false
		return
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&AppConfig)
	IsConfigured = true
}

// SaveConfig saves the configuration to gofiles.json
func SaveConfig(username, password string) error {
	AppConfig = types.ConfigFile{Username: username, Password: password, CreatedAt: time.Now()}
	file, err := os.Create(ConfigFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(AppConfig)
}

// GetEnv helper to get env variables
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
