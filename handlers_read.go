package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
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
		if filterExt != "" && strings.ToLower(filepath.Ext(f.Name())) != filterExt {
			continue
		}
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

	err := filepath.WalkDir(fullStartPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(fullStartPath, path)
		if relPath == "." {
			return nil
		}

		// --- NAME SEARCH ---
		if searchType == "name" {
			if strings.Contains(strings.ToLower(d.Name()), query) {
				info, _ := d.Info()
				results = append(results, FileInfo{
					Name:    relPath,
					Size:    info.Size(),
					IsDir:   d.IsDir(),
					ModTime: info.ModTime().Format(time.RFC3339),
					Type:    filepath.Ext(d.Name()),
				})
			}
		}

		// --- CONTENT SEARCH ---
		if searchType == "content" && !d.IsDir() {
			info, _ := d.Info()
			if info.Size() > 5*1024*1024 {
				return nil
			} // Skip > 5MB

			file, err := os.Open(path)
			if err == nil {
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
						break
					}
				}
				file.Close()
			}
		}

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