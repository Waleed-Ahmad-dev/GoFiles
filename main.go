package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Config: The root folder we are browsing
const RootFolder = "."

// FileInfo defines the JSON structure for our API
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime string `json:"mod_time"`
	Type    string `json:"type"`
}

func main() {
	// --- READ Operations ---
	http.HandleFunc("/api/files", handleListFiles)       // List directory
	http.HandleFunc("/api/download", handleDownloadFile) // Download/Preview

	// --- WRITE Operations (New!) ---
	http.HandleFunc("/api/upload", handleUploadFile)     // Upload a file
	http.HandleFunc("/api/mkdir", handleCreateDir)       // Create a folder
	http.HandleFunc("/api/delete", handleDelete)         // Delete file/folder

	fmt.Println("🚀 GoFiles Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ---------------- HANDLERS ----------------

func handleListFiles(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	fullPath := filepath.Join(RootFolder, reqPath)

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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	fullPath := filepath.Join(RootFolder, reqPath)

	if !isPathSafe(fullPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, fullPath)
}

// handleUploadFile accepts multipart/form-data uploads
func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return // Preflight CORS check often sends OPTIONS
	}

	// 1. Limit memory usage to 10MB for parsing (rest goes to disk temp)
	r.ParseMultipartForm(10 << 20)

	// 2. Get the target directory from query params
	targetDir := r.URL.Query().Get("path")
	fullDirPath := filepath.Join(RootFolder, targetDir)

	if !isPathSafe(fullDirPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// 3. Retrieve the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 4. Create the destination file
	// filepath.Base ensures the user can't send a filename like "../../virus.exe"
	dstPath := filepath.Join(fullDirPath, filepath.Base(handler.Filename))
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Error creating file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 5. Stream the bits (Copy from Upload -> Disk)
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Uploaded: %s\n", handler.Filename)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upload successful"))
}

func handleCreateDir(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	// We expect a JSON body like: {"path": "folder1", "name": "new_folder"}
	var req struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(RootFolder, req.Path, req.Name)

	if !isPathSafe(fullPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	if err := os.Mkdir(fullPath, 0755); err != nil {
		http.Error(w, "Could not create directory", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// We allow DELETE method or POST method (some strict firewalls block DELETE)
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		return
	}

	targetPath := r.URL.Query().Get("path")
	fullPath := filepath.Join(RootFolder, targetPath)

	if !isPathSafe(fullPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// RemoveAll deletes both files and empty/non-empty folders
	if err := os.RemoveAll(fullPath); err != nil {
		http.Error(w, "Could not delete item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ---------------- HELPERS ----------------

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func isPathSafe(path string) bool {
	root, err := filepath.Abs(RootFolder)
	if err != nil {
		return false
	}
	target, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	isOutside := len(rel) >= 2 && rel[0:2] == ".."
	return !isOutside
}