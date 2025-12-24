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

// CreateDirRequest defines the JSON body for creating a folder
type CreateDirRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

// ActionRequest handles Move, Copy, and Rename
// SourcePath: Where the file is now
// DestPath: Where you want it to go (for Copy/Move)
// NewName: The new name (for Rename)
type ActionRequest struct {
	SourcePath string `json:"sourcePath"`
	DestPath   string `json:"destPath"` 
	NewName    string `json:"newName"`
}

type TrashInfo struct {
	OriginalPath string    `json:"originalPath"`
	DeletedAt    time.Time `json:"deletedAt"`
	Filename     string    `json:"filename"` // The unique name in trash
}