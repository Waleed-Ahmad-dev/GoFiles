package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Register Routes (handlers are now in handlers.go)
	http.HandleFunc("/api/files", handleListFiles)
	http.HandleFunc("/api/download", handleDownloadFile)
	http.HandleFunc("/api/upload", handleUploadFile)
	http.HandleFunc("/api/mkdir", handleCreateDir)
	http.HandleFunc("/api/delete", handleDelete)

	fmt.Println("🚀 GoFiles Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}