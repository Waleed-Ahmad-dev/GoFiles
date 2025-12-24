package main

import (
	"encoding/json"
	"os"
	"time"
)

// Global Config
const RootFolder = "."
const TrashFolder = ".trash"
const TrashRetention = 30 * 24 * time.Hour
const ConfigFileName = "gofiles.json"

// Runtime State
var AppConfig ConfigFile
var IsConfigured = false

// InitConfig tries to load gofiles.json
func InitConfig() {
	file, err := os.Open(ConfigFileName)
	if err != nil {
		// File doesn't exist? That means we need Setup!
		IsConfigured = false
		return
	}
	defer file.Close()

	// File exists? Load credentials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		// If JSON is broken, force setup
		IsConfigured = false
		return
	}

	IsConfigured = true
}

// SaveConfig writes the credentials to disk
func SaveConfig(username, password string) error {
	AppConfig = ConfigFile{
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
	}

	file, err := os.Create(ConfigFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(AppConfig)
}