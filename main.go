package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize Trash (Creates folder + Starts Auto-Cleanup Timer)
	InitTrash()

	// Read
	http.HandleFunc("/api/files", handleListFiles)
	http.HandleFunc("/api/download", handleDownloadFile)
	http.HandleFunc("/api/trash/list", handleListTrash) // NEW

	// Write
	http.HandleFunc("/api/upload", handleUploadFile)
	http.HandleFunc("/api/mkdir", handleCreateDir)

	// Delete & Restore
	http.HandleFunc("/api/delete", handleDelete)          // Modified (Soft delete)
	http.HandleFunc("/api/trash/restore", handleRestore)  // NEW
	http.HandleFunc("/api/trash/empty", handleEmptyTrash) // NEW

	// Organize
	http.HandleFunc("/api/rename", handleRename)
	http.HandleFunc("/api/move", handleMove)
	http.HandleFunc("/api/copy", handleCopy)

	fmt.Println("🚀 GoFiles Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
