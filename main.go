package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize Trash
	InitTrash()

	// Read & Search
	http.HandleFunc("/api/files", handleListFiles)
	http.HandleFunc("/api/download", handleDownloadFile)
	http.HandleFunc("/api/search", handleSearch) // NEW: Search Endpoint

	// Trash
	http.HandleFunc("/api/trash/list", handleListTrash)
	http.HandleFunc("/api/trash/restore", handleRestore)
	http.HandleFunc("/api/trash/empty", handleEmptyTrash)

	// Write
	http.HandleFunc("/api/upload", handleUploadFile)
	http.HandleFunc("/api/mkdir", handleCreateDir)
	http.HandleFunc("/api/delete", handleDelete)

	// Organize
	http.HandleFunc("/api/rename", handleRename)
	http.HandleFunc("/api/move", handleMove)
	http.HandleFunc("/api/copy", handleCopy)

	fmt.Println("🚀 GoFiles Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}