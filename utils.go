package main

import (
	"net/http"
	"path/filepath"
)

// enableCors allows your future Frontend to talk to this backend
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// isPathSafe ensures the user doesn't try to access protected folders
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