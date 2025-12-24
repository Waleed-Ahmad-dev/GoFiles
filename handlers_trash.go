package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func handleListTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	trashRoot := filepath.Join(RootFolder, TrashFolder)

	files, _ := ioutil.ReadDir(trashRoot)

	var trashList []TrashInfo
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			metaBytes, _ := ioutil.ReadFile(filepath.Join(trashRoot, f.Name()))
			var meta TrashInfo
			json.Unmarshal(metaBytes, &meta)
			trashList = append(trashList, meta)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trashList)
}

func handleRestore(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	trashFilename := r.URL.Query().Get("name")

	// Basic security check
	if strings.Contains(trashFilename, "/") || strings.Contains(trashFilename, "\\") {
		return
	}

	RestoreFromTrash(trashFilename)
	w.WriteHeader(http.StatusOK)
}

func handleEmptyTrash(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	os.RemoveAll(filepath.Join(RootFolder, TrashFolder))
	os.Mkdir(filepath.Join(RootFolder, TrashFolder), 0755)
	w.WriteHeader(http.StatusOK)
}