package main

import (
	"encoding/json"
	"os"
	"time"
)

// Global Config
const RootFolder = "."
const TrashFolder = ".trash"
const ThumbsFolder = ".thumbs" // NEW: Hidden folder for thumbnails
const TrashRetention = 30 * 24 * time.Hour
const ConfigFileName = "gofiles.json"

// Runtime State
var AppConfig ConfigFile
var IsConfigured = false

// InitConfig tries to load gofiles.json
func InitConfig() {
	// ... (Keep existing logic) ...
	file, err := os.Open(ConfigFileName)
	if err != nil {
		IsConfigured = false
		return
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&AppConfig)
	IsConfigured = true
}

// ... (Keep SaveConfig and getEnv) ...
func SaveConfig(username, password string) error {
	AppConfig = ConfigFile{Username: username, Password: password, CreatedAt: time.Now()}
	file, err := os.Create(ConfigFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(AppConfig)
}

// Helper to get env variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}