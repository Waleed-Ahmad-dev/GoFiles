package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	r.ParseMultipartForm(10 << 20)
	targetDir := r.URL.Query().Get("path")
	fullDirPath := filepath.Join(RootFolder, targetDir)

	if !isPathSafe(fullDirPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		return
	}
	defer file.Close()

	dstPath := filepath.Join(fullDirPath, filepath.Base(handler.Filename))
	dst, err := os.Create(dstPath)
	if err != nil {
		return
	}
	defer dst.Close()

	io.Copy(dst, file)
	w.WriteHeader(http.StatusOK)
}

func handleCreateDir(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req CreateDirRequest
	json.NewDecoder(r.Body).Decode(&req)
	fullPath := filepath.Join(RootFolder, req.Path, req.Name)

	if !isPathSafe(fullPath) {
		return
	}
	os.Mkdir(fullPath, 0755)
	w.WriteHeader(http.StatusOK)
}