package main

import (
	"os"
	"time"
)

// Config
const RootFolder = "."
const TrashFolder = ".trash"
const TrashRetention = 30 * 24 * time.Hour

// Auth Config (Defaults, but should be changed via Environment Variables)
var AdminUser = getEnv("GOFILES_USER", "admin")
var AdminPass = getEnv("GOFILES_PASS", "admin123")

// Helper to get env variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}