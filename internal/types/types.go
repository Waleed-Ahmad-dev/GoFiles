package types

import "time"

// FileInfo represents the details of a file or directory
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime string `json:"mod_time"`
	Type    string `json:"type"`
}

// ConfigFile represents the structure of the configuration file
type ConfigFile struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateDirRequest represents the request body for creating a directory
type CreateDirRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

// ActionRequest represents a generic file action (rename, move, copy)
type ActionRequest struct {
	SourcePath string `json:"sourcePath"`
	DestPath   string `json:"destPath"`
	NewName    string `json:"newName"`
}

// TrashInfo represents metadata for a trashed file
type TrashInfo struct {
	OriginalPath string    `json:"originalPath"`
	DeletedAt    time.Time `json:"deletedAt"`
	Filename     string    `json:"filename"`
}

// ArchiveRequest represents the request for Zip/Unzip operations
type ArchiveRequest struct {
	SourcePath string `json:"sourcePath"` // File/Folder to zip, or Zip file to unzip
	DestPath   string `json:"destPath"`   // Where to save
	Password   string `json:"password"`   // Optional: Leave empty for no password
}

// SaveFileRequest represents the request to save a text file
type SaveFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// LoginRequest represents the login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
