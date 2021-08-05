package main

import (
	"github.com/alvanrahimli/dots-server/models"
	"github.com/gorilla/mux"
	"net/http"
	"path"
)

func getArchiveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	archiveName := vars["name"]

	InfoLogger.Printf("File '%s' requested", archiveName)
	http.ServeFile(w, r, path.Join(models.ArchivesFolderRoot, archiveName))
	InfoLogger.Printf("File '%s' served", archiveName)
}
