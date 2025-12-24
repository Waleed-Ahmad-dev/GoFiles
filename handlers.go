package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// handleListFiles displays files in a folder, with optional filtering
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

	// --- FILTERING PARAMETERS ---
	filterExt := strings.ToLower(r.URL.Query().Get("ext")) // e.g. ".jpg"
	minSizeStr := r.URL.Query().Get("min_size")            // e.g. "1048576" (1MB)

	var minSize int64 = 0
	if minSizeStr != "" {
		minSize, _ = strconv.ParseInt(minSizeStr, 10, 64)
	}

	files, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusNotFound)
		return
	}

	var fileList []FileInfo
	for _, f := range files {
		info, _ := f.Info()

		// --- APPLY FILTERS ---
		// 1. Extension Filter
		if filterExt != "" && strings.ToLower(filepath.Ext(f.Name())) != filterExt {
			continue
		}
		// 2. Size Filter
		if minSize > 0 && info.Size() < minSize {
			continue
		}

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

// handleSearch performs recursive search for Name or Content
func handleSearch(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodGet {
		return
	}

	query := strings.ToLower(r.URL.Query().Get("q"))
	searchType := r.URL.Query().Get("type") // "name" or "content"
	startPath := r.URL.Query().Get("path")

	if query == "" {
		http.Error(w, "Query is empty", http.StatusBadRequest)
		return
	}

	fullStartPath := filepath.Join(RootFolder, startPath)
	if !isPathSafe(fullStartPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	var results []FileInfo

	// filepath.WalkDir is efficient and recursive
	err := filepath.WalkDir(fullStartPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		} // Skip permission errors

		// Calculate relative path for display (e.g. "subfolder/image.jpg")
		relPath, _ := filepath.Rel(fullStartPath, path)
		if relPath == "." {
			return nil
		}

		// --- NAME SEARCH ---
		if searchType == "name" {
			if strings.Contains(strings.ToLower(d.Name()), query) {
				info, _ := d.Info()
				results = append(results, FileInfo{
					Name:    relPath, // Return relative path so UI knows where it is
					Size:    info.Size(),
					IsDir:   d.IsDir(),
					ModTime: info.ModTime().Format(time.RFC3339),
					Type:    filepath.Ext(d.Name()),
				})
			}
		}

		// --- CONTENT SEARCH ---
		if searchType == "content" && !d.IsDir() {
			// Optimization: Skip files > 5MB to avoid freezing the server
			info, _ := d.Info()
			if info.Size() > 5*1024*1024 {
				return nil
			}

			// Read file content
			file, err := os.Open(path)
			if err == nil {
				// Use Scanner to check line by line (memory efficient)
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					if strings.Contains(strings.ToLower(scanner.Text()), query) {
						results = append(results, FileInfo{
							Name:    relPath,
							Size:    info.Size(),
							IsDir:   false,
							ModTime: info.ModTime().Format(time.RFC3339),
							Type:    filepath.Ext(d.Name()),
						})
						break // Found match, stop reading this file
					}
				}
				file.Close()
			}
		}

		// Stop if we found too many results (Safety limit)
		if len(results) > 100 {
			return io.EOF
		}

		return nil
	})

	if err != nil && err != io.EOF {
		fmt.Println("Search error:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodGet {
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

// ... (Keep handleRestore, handleListTrash, handleEmptyTrash, handleRename, handleMove, handleCopy exactly as they were) ...
// (I omitted them here to save space, but make sure you keep them in the file!)
func handleRestore(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}
	trashFilename := r.URL.Query().Get("name")
	if strings.Contains(trashFilename, "/") {
		return
	}
	RestoreFromTrash(trashFilename)
	w.WriteHeader(http.StatusOK)
}
func handleListTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	trashRoot := filepath.Join(RootFolder, TrashFolder)
	files, _ := ioutil.ReadDir(trashRoot)
	var trashList []TrashInfo
	for _, f := range files {
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
func handleEmptyTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}
	os.RemoveAll(filepath.Join(RootFolder, TrashFolder))
	os.Mkdir(filepath.Join(RootFolder, TrashFolder), 0755)
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