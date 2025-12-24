package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// 1. Initialize Sub-systems
	InitTrash()
	InitConfig()

	// Ensure Thumbs folder exists
	os.MkdirAll(filepath.Join(RootFolder, ThumbsFolder), 0755)

	// --- PUBLIC ROUTES ---
	http.HandleFunc("/api/system/status", handleSystemStatus)
	http.HandleFunc("/api/setup", handleSetup)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/logout", handleLogout)

	// --- PROTECTED ROUTES ---
	http.HandleFunc("/api/me", AuthMiddleware(handleCheckAuth))

	// Media
	http.HandleFunc("/api/thumbnail", AuthMiddleware(handleThumbnail))

	// Read & Search
	http.HandleFunc("/api/files", AuthMiddleware(handleListFiles))
	http.HandleFunc("/api/download", AuthMiddleware(handleDownloadFile))
	http.HandleFunc("/api/search", AuthMiddleware(handleSearch))

	// Zip / Unzip (NEW)
	http.HandleFunc("/api/zip", AuthMiddleware(handleZip))
	http.HandleFunc("/api/unzip", AuthMiddleware(handleUnzip))

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
		fmt.Println("⚠️  SYSTEM NOT CONFIGURED. Go to http://localhost:8080 to set up.")
	} else {
		fmt.Println("✅ System configured.")
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}