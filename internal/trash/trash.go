package trash

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"GoFiles/internal/config"
	"GoFiles/internal/types"
)

// InitTrash creates the hidden trash folder if it doesn't exist
// and starts the background cleanup task.
func InitTrash() {
	trashPath := filepath.Join(config.RootFolder, config.TrashFolder)
	if _, err := os.Stat(trashPath); os.IsNotExist(err) {
		os.Mkdir(trashPath, 0755)
	}

	// Start background cleanup (Runs in a separate thread)
	go startTrashCleanup()
}

// MoveToTrash performs a "Soft Delete"
func MoveToTrash(relativePath string) error {
	fullSourcePath := filepath.Join(config.RootFolder, relativePath)
	trashRoot := filepath.Join(config.RootFolder, config.TrashFolder)

	// 1. Generate unique name (file.txt -> file.txt_1739281)
	info, err := os.Stat(fullSourcePath)
	if err != nil {
		return err
	}
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	trashName := info.Name() + "_" + timestamp
	trashPath := filepath.Join(trashRoot, trashName)

	// 2. Create Metadata File (.json)
	meta := types.TrashInfo{
		OriginalPath: relativePath,
		DeletedAt:    time.Now(),
		Filename:     trashName,
	}
	metaBytes, _ := json.MarshalIndent(meta, "", "  ")

	// Save metadata: .trash/file.txt_1739281.json
	err = ioutil.WriteFile(trashPath+".json", metaBytes, 0644)
	if err != nil {
		return err
	}

	// 3. Move the actual file
	return os.Rename(fullSourcePath, trashPath)
}

// RestoreFromTrash moves a file back to its original location
func RestoreFromTrash(trashFilename string) error {
	trashRoot := filepath.Join(config.RootFolder, config.TrashFolder)
	trashFilePath := filepath.Join(trashRoot, trashFilename)
	metaFilePath := trashFilePath + ".json"

	// 1. Read Metadata
	metaBytes, err := ioutil.ReadFile(metaFilePath)
	if err != nil {
		return fmt.Errorf("metadata not found")
	}
	var meta types.TrashInfo
	json.Unmarshal(metaBytes, &meta)

	// 2. Check if original folder still exists
	destPath := filepath.Join(config.RootFolder, meta.OriginalPath)
	destDir := filepath.Dir(destPath)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		// If original folder is gone, recreate it
		os.MkdirAll(destDir, 0755)
	}

	// 3. Move File Back
	if err := os.Rename(trashFilePath, destPath); err != nil {
		return err
	}

	// 4. Delete Metadata File
	os.Remove(metaFilePath)
	return nil
}

// startTrashCleanup runs forever, checking for old files every hour
func startTrashCleanup() {
	for {
		// Sleep first to let server start up
		time.Sleep(1 * time.Hour)

		fmt.Println("🧹 Running Auto-Trash Cleanup...")
		trashRoot := filepath.Join(config.RootFolder, config.TrashFolder)
		files, _ := ioutil.ReadDir(trashRoot)

		for _, f := range files {
			// specific logic: only check .json files to find age
			if strings.HasSuffix(f.Name(), ".json") {
				continue
			}

			// Check the corresponding JSON file for the date
			metaBytes, err := ioutil.ReadFile(filepath.Join(trashRoot, f.Name()+".json"))
			if err != nil {
				// No metadata? Just rely on file mod time
				if time.Since(f.ModTime()) > config.TrashRetention {
					os.RemoveAll(filepath.Join(trashRoot, f.Name()))
				}
				continue
			}

			var meta types.TrashInfo
			json.Unmarshal(metaBytes, &meta)

			if time.Since(meta.DeletedAt) > config.TrashRetention {
				fmt.Printf("🗑️ Auto-deleting old file: %s\n", f.Name())
				// Delete File AND Metadata
				os.RemoveAll(filepath.Join(trashRoot, f.Name()))
				os.Remove(filepath.Join(trashRoot, f.Name()+".json"))
			}
		}
	}
}
