package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"GoFiles/internal/config"
	"GoFiles/internal/trash"
	"GoFiles/internal/types"
	"GoFiles/internal/utils"
)

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		return
	}

	targetPath := r.URL.Query().Get("path")
	permanent := r.URL.Query().Get("permanent") == "true"

	if !utils.IsPathSafe(filepath.Join(config.RootFolder, targetPath)) {
		return
	}

	if permanent {
		os.RemoveAll(filepath.Join(config.RootFolder, targetPath))
	} else {
		trash.MoveToTrash(targetPath)
	}
	w.WriteHeader(http.StatusOK)
}

func HandleRename(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req types.ActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	oldPath := filepath.Join(config.RootFolder, req.SourcePath)
	newPath := filepath.Join(filepath.Dir(oldPath), req.NewName)

	if !utils.IsPathSafe(oldPath) || !utils.IsPathSafe(newPath) {
		return
	}

	os.Rename(oldPath, newPath)
	w.WriteHeader(http.StatusOK)
}

func HandleMove(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req types.ActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	srcPath := filepath.Join(config.RootFolder, req.SourcePath)
	destPath := filepath.Join(config.RootFolder, req.DestPath, filepath.Base(req.SourcePath))

	if !utils.IsPathSafe(srcPath) || !utils.IsPathSafe(destPath) {
		return
	}

	os.Rename(srcPath, destPath)
	w.WriteHeader(http.StatusOK)
}

func HandleCopy(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req types.ActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	srcPath := filepath.Join(config.RootFolder, req.SourcePath)
	destPath := filepath.Join(config.RootFolder, req.DestPath, filepath.Base(req.SourcePath))

	if !utils.IsPathSafe(srcPath) || !utils.IsPathSafe(destPath) {
		return
	}

	info, _ := os.Stat(srcPath)
	if info.IsDir() {
		utils.CopyDir(srcPath, destPath)
	} else {
		utils.CopyFile(srcPath, destPath)
	}
	w.WriteHeader(http.StatusOK)
}
