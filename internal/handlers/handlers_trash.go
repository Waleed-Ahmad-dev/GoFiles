package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"GoFiles/internal/config"
	"GoFiles/internal/trash"
	"GoFiles/internal/types"
	"GoFiles/internal/utils"
)

func HandleListTrash(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	trashRoot := filepath.Join(config.RootFolder, config.TrashFolder)

	files, _ := ioutil.ReadDir(trashRoot)

	var trashList []types.TrashInfo
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			metaBytes, _ := ioutil.ReadFile(filepath.Join(trashRoot, f.Name()))
			var meta types.TrashInfo
			json.Unmarshal(metaBytes, &meta)
			trashList = append(trashList, meta)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trashList)
}

func HandleRestore(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	trashFilename := r.URL.Query().Get("name")

	// Basic security check
	if strings.Contains(trashFilename, "/") || strings.Contains(trashFilename, "\\") {
		return
	}

	trash.RestoreFromTrash(trashFilename)
	w.WriteHeader(http.StatusOK)
}

func HandleEmptyTrash(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	os.RemoveAll(filepath.Join(config.RootFolder, config.TrashFolder))
	os.Mkdir(filepath.Join(config.RootFolder, config.TrashFolder), 0755)
	w.WriteHeader(http.StatusOK)
}
