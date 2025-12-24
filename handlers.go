package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

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

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB limit

	targetDir := r.URL.Query().Get("path")
	fullDirPath := filepath.Join(RootFolder, targetDir)

	if !isPathSafe(fullDirPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dstPath := filepath.Join(fullDirPath, filepath.Base(handler.Filename))
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Error creating file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

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

	var req CreateDirRequest
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
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		return
	}

	targetPath := r.URL.Query().Get("path")
	fullPath := filepath.Join(RootFolder, targetPath)

	if !isPathSafe(fullPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	if err := os.RemoveAll(fullPath); err != nil {
		http.Error(w, "Could not delete item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}