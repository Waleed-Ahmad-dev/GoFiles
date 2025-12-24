package main

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