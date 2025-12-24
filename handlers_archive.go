package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yeka/zip" // Replaces standard archive/zip
)

// handleZip compresses a file or folder into a .zip
func handleZip(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req ArchiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	srcPath := filepath.Join(RootFolder, req.SourcePath)
	// Destination: If DestPath is empty, save next to source
	destPath := ""
	if req.DestPath != "" {
		destPath = filepath.Join(RootFolder, req.DestPath)
	} else {
		destPath = srcPath + ".zip"
	}

	if !isPathSafe(srcPath) || !isPathSafe(destPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// Create the Zip File
	zipFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Could not create zip file", http.StatusInternalServerError)
		return
	}
	defer zipFile.Close()

	// Initialize Zip Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the source directory/file
	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Don't zip the zip file itself if it's in the same folder
		if path == destPath {
			return nil
		}

		// Create header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Make path relative to the root of the archive
		// e.g. zipping /users/docs/work -> work/resume.pdf
		relPath, _ := filepath.Rel(filepath.Dir(srcPath), path)
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate // Compress
		}

		// Set Password if provided
		if req.Password != "" {
			header.SetPassword(req.Password)
		}

		// Create writer for this file inside zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Copy content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		http.Error(w, "Error zipping: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleUnzip extracts a zip file
func handleUnzip(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req ArchiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	srcPath := filepath.Join(RootFolder, req.SourcePath)
	destPath := filepath.Join(RootFolder, req.DestPath)

	if !isPathSafe(srcPath) || !isPathSafe(destPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// Open Zip Reader
	reader, err := zip.OpenReader(srcPath)
	if err != nil {
		http.Error(w, "Failed to open zip: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer reader.Close()

	// Iterate through files in zip
	for _, file := range reader.File {
		// Set Password if needed
		if file.IsEncrypted() {
			file.SetPassword(req.Password)
		}

		// Calculate extract path
		fpath := filepath.Join(destPath, file.Name)

		// Zip Slip Protection (Security)
		// Prevent zips from containing "../../virus.exe"
		if !strings.HasPrefix(fpath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			continue // Skip illegal paths
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make parent dirs
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			http.Error(w, "File permission error", http.StatusInternalServerError)
			return
		}

		// Open file inside zip
		rc, err := file.Open()
		if err != nil {
			if strings.Contains(err.Error(), "password") {
				http.Error(w, "Incorrect Password", http.StatusUnauthorized)
			} else {
				http.Error(w, "Read error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Create file on disk
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			return
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()
		if err != nil {
			http.Error(w, "Extract error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
