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
	if r.Method != http.MethodDelete && r.Method != http.MethodPost { return }

	targetPath := r.URL.Query().Get("path")
	permanent := r.URL.Query().Get("permanent") == "true"

	// Security Check
	if !isPathSafe(filepath.Join(RootFolder, targetPath)) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// If user specifically asked for permanent delete (Shift+Delete style)
	if permanent {
		fullPath := filepath.Join(RootFolder, targetPath)
		if err := os.RemoveAll(fullPath); err != nil {
			http.Error(w, "Could not delete", http.StatusInternalServerError)
			return
		}
	} else {
		// Normal Delete -> Send to Trash
		if err := MoveToTrash(targetPath); err != nil {
			http.Error(w, "Could not move to trash: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// ---------------- NEW HANDLERS ----------------

// handleRestore restores a file given its ID (trash filename)
func handleRestore(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost { return }

	trashFilename := r.URL.Query().Get("name")
	
	// Basic security: ensure we are only touching files in .trash
	if strings.Contains(trashFilename, "/") || strings.Contains(trashFilename, "\\") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	if err := RestoreFromTrash(trashFilename); err != nil {
		http.Error(w, "Restore failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleListTrash shows what is in the bin
func handleListTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	trashRoot := filepath.Join(RootFolder, TrashFolder)
	
	files, _ := ioutil.ReadDir(trashRoot)
	
	var trashList []TrashInfo

	for _, f := range files {
		// We only care about the .json files to build our list
		if strings.HasSuffix(f.Name(), ".json") {
			metaBytes, _ := ioutil.ReadFile(filepath.Join(trashRoot, f.Name()))
			var meta TrashInfo
			json.Unmarshal(metaBytes, &meta)
			trashList = append(trashList, meta)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trashList)
}

// handleEmptyTrash permanently deletes EVERYTHING in .trash
func handleEmptyTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost && r.Method != http.MethodDelete { return }

	trashRoot := filepath.Join(RootFolder, TrashFolder)
	
	// Delete the whole folder and recreate it
	os.RemoveAll(trashRoot)
	os.Mkdir(trashRoot, 0755)

	w.WriteHeader(http.StatusOK)
}

func handleRename(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost { return }

	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	oldPath := filepath.Join(RootFolder, req.SourcePath)
	// New path is just the directory of the old path + the new name
	newPath := filepath.Join(filepath.Dir(oldPath), req.NewName)

	if !isPathSafe(oldPath) || !isPathSafe(newPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		http.Error(w, "Could not rename file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost { return }

	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	srcPath := filepath.Join(RootFolder, req.SourcePath)
	destPath := filepath.Join(RootFolder, req.DestPath, filepath.Base(req.SourcePath))

	if !isPathSafe(srcPath) || !isPathSafe(destPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// os.Rename moves files (and is very fast)
	if err := os.Rename(srcPath, destPath); err != nil {
		http.Error(w, "Could not move file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleCopy(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost { return }

	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	srcPath := filepath.Join(RootFolder, req.SourcePath)
	destPath := filepath.Join(RootFolder, req.DestPath, filepath.Base(req.SourcePath))

	if !isPathSafe(srcPath) || !isPathSafe(destPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// Check if source is a file or directory
	info, err := os.Stat(srcPath)
	if err != nil {
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	if info.IsDir() {
		err = CopyDir(srcPath, destPath)
	} else {
		err = CopyFile(srcPath, destPath)
	}

	if err != nil {
		http.Error(w, "Error copying: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}