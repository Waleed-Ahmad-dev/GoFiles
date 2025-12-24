package main

import "time"

// FileInfo defines the JSON structure for our API
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime string `json:"mod_time"`
	Type    string `json:"type"`
}

// ConfigFile is what we save to disk (gofiles.json)
type ConfigFile struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"` // In production, we should hash this!
	CreatedAt time.Time `json:"created_at"`
}

// ... (Keep CreateDirRequest, ActionRequest, TrashInfo as they were) ...
type CreateDirRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type ActionRequest struct {
	SourcePath string `json:"sourcePath"`
	DestPath   string `json:"destPath"`
	NewName    string `json:"newName"`
}

type TrashInfo struct {
	OriginalPath string    `json:"originalPath"`
	DeletedAt    time.Time `json:"deletedAt"`
	Filename     string    `json:"filename"`
}