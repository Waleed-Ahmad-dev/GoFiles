package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 1. Initialize Sub-systems
	InitTrash()
	InitConfig() // Checks if gofiles.json exists

	// --- PUBLIC ROUTES (No Auth Required) ---
	http.HandleFunc("/api/system/status", handleSystemStatus) // Checks if we need setup
	http.HandleFunc("/api/setup", handleSetup)                // Performs the setup
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/logout", handleLogout)

	// --- PROTECTED ROUTES ---
	http.HandleFunc("/api/me", AuthMiddleware(handleCheckAuth))

	// Read & Search
	http.HandleFunc("/api/files", AuthMiddleware(handleListFiles))
	http.HandleFunc("/api/download", AuthMiddleware(handleDownloadFile))
	http.HandleFunc("/api/search", AuthMiddleware(handleSearch))

	// Trash
	http.HandleFunc("/api/trash/list", AuthMiddleware(handleListTrash))
	http.HandleFunc("/api/trash/restore", AuthMiddleware(handleRestore))
	http.HandleFunc("/api/trash/empty", AuthMiddleware(handleEmptyTrash))

	// Write
	http.HandleFunc("/api/upload", AuthMiddleware(handleUploadFile))
	http.HandleFunc("/api/mkdir", AuthMiddleware(handleCreateDir))
	http.HandleFunc("/api/delete", AuthMiddleware(handleDelete))

	// Organize
	http.HandleFunc("/api/rename", AuthMiddleware(handleRename))
	http.HandleFunc("/api/move", AuthMiddleware(handleMove))
	http.HandleFunc("/api/copy", AuthMiddleware(handleCopy))

	fmt.Println("🚀 GoFiles Server started on http://localhost:8080")
	if !IsConfigured {
		fmt.Println("⚠️  SYSTEM NOT CONFIGURED. Please go to the UI to set up an admin account.")
	} else {
		fmt.Println("✅ System configured. Login enabled.")
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}