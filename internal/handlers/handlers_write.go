package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"GoFiles/internal/config"
	"GoFiles/internal/types"
	"GoFiles/internal/utils"
)

func HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	r.ParseMultipartForm(10 << 20)
	targetDir := r.URL.Query().Get("path")
	fullDirPath := filepath.Join(config.RootFolder, targetDir)

	if !utils.IsPathSafe(fullDirPath) {
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

func HandleCreateDir(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req types.CreateDirRequest
	json.NewDecoder(r.Body).Decode(&req)
	fullPath := filepath.Join(config.RootFolder, req.Path, req.Name)

	if !utils.IsPathSafe(fullPath) {
		return
	}
	os.Mkdir(fullPath, 0755)
	w.WriteHeader(http.StatusOK)
}

func HandleSaveFile(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req types.SaveFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(config.RootFolder, req.Path)

	if !utils.IsPathSafe(fullPath) {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// Write the string content to the file
	err := ioutil.WriteFile(fullPath, []byte(req.Content), 0644)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
