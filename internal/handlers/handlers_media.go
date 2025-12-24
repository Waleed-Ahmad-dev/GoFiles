package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"GoFiles/internal/config"
	"GoFiles/internal/utils"

	"github.com/disintegration/imaging"
)

// HandleThumbnail generates or retrieves a cached thumbnail
func HandleThumbnail(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodGet {
		return
	}

	reqPath := r.URL.Query().Get("path")
	fullPath := filepath.Join(config.RootFolder, reqPath)

	if !utils.IsPathSafe(fullPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// 1. Check if the file is actually an image
	ext := strings.ToLower(filepath.Ext(fullPath))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		http.Error(w, "Not an image", http.StatusBadRequest)
		return
	}

	// 2. Get File Info (to check modification time)
	info, err := os.Stat(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// 3. Generate a Unique Cache Filename
	// We hash the Path + ModTime. If the file is edited, ModTime changes, hash changes -> New Thumbnail!
	hashKey := fmt.Sprintf("%s-%d", fullPath, info.ModTime().Unix())
	hasher := md5.New()
	hasher.Write([]byte(hashKey))
	hash := hex.EncodeToString(hasher.Sum(nil))

	thumbFilename := hash + ".jpg"
	thumbPath := filepath.Join(config.RootFolder, config.ThumbsFolder, thumbFilename)

	// 4. Check if Thumbnail already exists on disk
	if _, err := os.Stat(thumbPath); err == nil {
		// HIT! Serve directly from cache
		http.ServeFile(w, r, thumbPath)
		return
	}

	// 5. MISS! Generate it.
	// Ensure .thumbs folder exists
	os.MkdirAll(filepath.Join(config.RootFolder, config.ThumbsFolder), 0755)

	// Open and Resize
	// imaging.Open handles rotation automatically (EXIF data)
	srcImage, err := imaging.Open(fullPath)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return
	}

	// Resize to width 300px, preserve aspect ratio
	// Lanczos is the best quality filter
	dstImage := imaging.Resize(srcImage, 300, 0, imaging.Lanczos)

	// Save to .thumbs folder
	err = imaging.Save(dstImage, thumbPath)
	if err != nil {
		http.Error(w, "Failed to save thumbnail", http.StatusInternalServerError)
		return
	}

	// Serve the newly created file
	http.ServeFile(w, r, thumbPath)
}
