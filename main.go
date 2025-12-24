package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Config: The root folder we are browsing
// "." means the current folder where the program is running
const RootFolder = "."

// FileInfo defines the JSON structure for our API
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime string `json:"mod_time"`
	Type    string `json:"type"` // Added file extension/type
}

func main() {
	// 1. List Files Endpoint
	http.HandleFunc("/api/files", handleListFiles)
	
	// 2. Download/View File Endpoint
	http.HandleFunc("/api/download", handleDownloadFile)

	fmt.Println("🚀 GoFiles Server started on http://localhost:8080")
	// Listen on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ---------------- HANDLERS ----------------

func handleListFiles(w http.ResponseWriter, r *http.Request) {
	enableCors(&w) // Enable frontend access

	// Get path from query, e.g., ?path=folder1
	reqPath := r.URL.Query().Get("path")
	fullPath := filepath.Join(RootFolder, reqPath)

	// Basic Security: Prevent going up directories
	if !isPathSafe(fullPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	files, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusNotFound)
		return
	}

	var fileList []FileInfo
	for _, f := range files {
		info, _ := f.Info()
		fileList = append(fileList, FileInfo{
			Name:    f.Name(),
			Size:    info.Size(),
			IsDir:   f.IsDir(),
			ModTime: info.ModTime().Format(time.RFC3339),
			Type:    filepath.Ext(f.Name()), 
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileList)
}

func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// Get file path
	reqPath := r.URL.Query().Get("path")
	fullPath := filepath.Join(RootFolder, reqPath)

	// Security Check
	if !isPathSafe(fullPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// Serve the file!
	// This handles streaming, ranges, and content-types automatically.
	http.ServeFile(w, r, fullPath)
}

// ---------------- HELPERS ----------------

// enableCors allows your future Frontend to talk to this backend
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// isPathSafe ensures the user doesn't try to access protected folders
func isPathSafe(path string) bool {
	// Get the absolute path of the root directory
	root, err := filepath.Abs(RootFolder)
	if err != nil {
		return false
	}

	// Get the absolute path of the requested file
	target, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Calculate the relative path from root to target
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}

	// If the relative path starts with "..", it means it's outside our root folder!
	// We also check if it is the root itself (".")
	isOutside := len(rel) >= 2 && rel[0:2] == ".." 
	
	return !isOutside
}