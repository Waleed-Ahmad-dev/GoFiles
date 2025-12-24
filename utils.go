package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// enableCors allows your future Frontend to talk to this backend
func enableCors(w *http.ResponseWriter) {
	// 1. Allow the specific origin sending the request (Dynamic Origin)
	// This is required when withCredentials is set to true
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	// 2. Allow credentials (cookies)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")

	// 3. Allowed Methods
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

	// 4. Allowed Headers
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// CopyDir recursively copies a directory tree
func CopyDir(src, dst string) error {
	// Get properties of source dir
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}
