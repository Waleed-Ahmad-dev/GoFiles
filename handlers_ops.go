package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

func handleDelete(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		return
	}

	targetPath := r.URL.Query().Get("path")
	permanent := r.URL.Query().Get("permanent") == "true"

	if !isPathSafe(filepath.Join(RootFolder, targetPath)) {
		return
	}

	if permanent {
		os.RemoveAll(filepath.Join(RootFolder, targetPath))
	} else {
		MoveToTrash(targetPath)
	}
	w.WriteHeader(http.StatusOK)
}

func handleRename(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req ActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	oldPath := filepath.Join(RootFolder, req.SourcePath)
	newPath := filepath.Join(filepath.Dir(oldPath), req.NewName)

	if !isPathSafe(oldPath) || !isPathSafe(newPath) {
		return
	}

	os.Rename(oldPath, newPath)
	w.WriteHeader(http.StatusOK)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req ActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	srcPath := filepath.Join(RootFolder, req.SourcePath)
	destPath := filepath.Join(RootFolder, req.DestPath, filepath.Base(req.SourcePath))

	if !isPathSafe(srcPath) || !isPathSafe(destPath) {
		return
	}

	os.Rename(srcPath, destPath)
	w.WriteHeader(http.StatusOK)
}

func handleCopy(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req ActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	srcPath := filepath.Join(RootFolder, req.SourcePath)
	destPath := filepath.Join(RootFolder, req.DestPath, filepath.Base(req.SourcePath))

	if !isPathSafe(srcPath) || !isPathSafe(destPath) {
		return
	}

	info, _ := os.Stat(srcPath)
	if info.IsDir() {
		CopyDir(srcPath, destPath)
	} else {
		CopyFile(srcPath, destPath)
	}
	w.WriteHeader(http.StatusOK)
}