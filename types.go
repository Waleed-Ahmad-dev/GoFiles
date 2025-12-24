package main

import "time"

// ... (Keep existing FileInfo, ConfigFile, CreateDirRequest) ...

type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime string `json:"mod_time"`
	Type    string `json:"type"`
}

type ConfigFile struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateDirRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

// ... (Keep ActionRequest, TrashInfo) ...

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

// NEW: For Zip/Unzip operations
type ArchiveRequest struct {
	SourcePath string `json:"sourcePath"` // File/Folder to zip, or Zip file to unzip
	DestPath   string `json:"destPath"`   // Where to save
	Password   string `json:"password"`   // Optional: Leave empty for no password
}
