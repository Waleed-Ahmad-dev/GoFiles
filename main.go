package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"GoFiles/internal/auth"
	"GoFiles/internal/config"
	"GoFiles/internal/handlers"
	"GoFiles/internal/trash"
)

func main() {
	// 1. Initialize Sub-systems
	trash.InitTrash()
	config.InitConfig()

	// Ensure Thumbs folder exists
	os.MkdirAll(filepath.Join(config.RootFolder, config.ThumbsFolder), 0755)

	// --- PUBLIC ROUTES ---
	http.HandleFunc("/api/system/status", auth.HandleSystemStatus)
	http.HandleFunc("/api/setup", auth.HandleSetup)
	http.HandleFunc("/api/login", auth.HandleLogin)
	http.HandleFunc("/api/logout", auth.HandleLogout)

	// --- PROTECTED ROUTES ---
	http.HandleFunc("/api/me", auth.AuthMiddleware(auth.HandleCheckAuth))

	// Media
	http.HandleFunc("/api/thumbnail", auth.AuthMiddleware(handlers.HandleThumbnail))

	// Read & Search
	http.HandleFunc("/api/files", auth.AuthMiddleware(handlers.HandleListFiles))
	http.HandleFunc("/api/download", auth.AuthMiddleware(handlers.HandleDownloadFile))
	http.HandleFunc("/api/download-zip", auth.AuthMiddleware(handlers.HandleDownloadZip)) // NEW: Stream Zip
	http.HandleFunc("/api/search", auth.AuthMiddleware(handlers.HandleSearch))

	// Zip / Unzip
	http.HandleFunc("/api/zip", auth.AuthMiddleware(handlers.HandleZip))
	http.HandleFunc("/api/unzip", auth.AuthMiddleware(handlers.HandleUnzip))

	// Trash
	http.HandleFunc("/api/trash/list", auth.AuthMiddleware(handlers.HandleListTrash))
	http.HandleFunc("/api/trash/restore", auth.AuthMiddleware(handlers.HandleRestore))
	http.HandleFunc("/api/trash/empty", auth.AuthMiddleware(handlers.HandleEmptyTrash))

	// Write
	http.HandleFunc("/api/upload", auth.AuthMiddleware(handlers.HandleUploadFile))
	http.HandleFunc("/api/save", auth.AuthMiddleware(handlers.HandleSaveFile)) // NEW: Text Save
	http.HandleFunc("/api/mkdir", auth.AuthMiddleware(handlers.HandleCreateDir))
	http.HandleFunc("/api/delete", auth.AuthMiddleware(handlers.HandleDelete))

	// Organize
	http.HandleFunc("/api/rename", auth.AuthMiddleware(handlers.HandleRename))
	http.HandleFunc("/api/move", auth.AuthMiddleware(handlers.HandleMove))
	http.HandleFunc("/api/copy", auth.AuthMiddleware(handlers.HandleCopy))

	fmt.Println("🚀 GoFiles Server started on http://localhost:8080")
	if !config.IsConfigured {
		fmt.Println("⚠️  SYSTEM NOT CONFIGURED. Go to http://localhost:8080 to set up.")
	} else {
		fmt.Println("✅ System configured.")
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}
