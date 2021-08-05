package main

import (
	"encoding/json"
	"fmt"
	"github.com/alvanrahimli/dots-server/dataaccess"
	"github.com/alvanrahimli/dots-server/models"
	"github.com/alvanrahimli/dots-server/utils"
	"io"
	"net/http"
	"os"
	"path"
)

func addPackageHandler(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Printf("Request received for URL %s", r.URL)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		_, err := fmt.Fprintf(w, "ParseForm() err: %v", err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ErrLogger.Println(err.Error())
			return
		}
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	packageName := r.FormValue("name")
	packageVersion := r.FormValue("version")
	file, fileHeader, fileErr := r.FormFile("archive")
	if fileErr != nil {
		ErrLogger.Println(fileErr.Error())
		http.Error(w, "Could not read archive file", http.StatusBadRequest)
		return
	}

	db := getDbInstance()
	defer db.Close()

	// Validate token
	userId, isValid := utils.ValidateToken(token, db)
	if !isValid {
		WarnLogger.Println("User '%s' token validation failed")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	archiveName := path.Join(models.ArchivesFolderRoot, path.Base(fileHeader.Filename))
	localFile, createErr := os.Create(archiveName)
	if createErr != nil {
		ErrLogger.Println(createErr)
		http.Error(w, "Could not save archive file", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()
	if _, copyErr := io.Copy(localFile, file); copyErr != nil {
		ErrLogger.Println(copyErr.Error())
		http.Error(w, "Could not save archive file", http.StatusInternalServerError)
		return
	}

	pack := models.Package{
		Id:          0,
		Name:        packageName,
		Version:     packageVersion,
		ArchiveName: archiveName,
		UserId:      userId,
	}

	response := models.HttpResponse{}

	validationErrors, isValid := utils.ValidatePackage(&pack)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		response = models.HttpResponse{
			Code:    1,
			Message: "One or more validation failed",
			Data: map[string]interface{}{
				"result": validationErrors,
			},
		}
	} else {
		// Add package
		pack, addErr := dataaccess.AddPackage(&pack, db)
		if addErr != nil {
			ErrLogger.Println(addErr)
			http.Error(w, "Could not add package", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		response = models.HttpResponse{
			Code:    0,
			Message: "Successfully created package",
			Data: map[string]interface{}{
				"package": pack,
			},
		}
	}

	responseJson, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		ErrLogger.Println(jsonErr.Error())
		http.Error(w, "Could not marshal response dto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, writeErr := w.Write(responseJson)
	if writeErr != nil {
		ErrLogger.Println(writeErr.Error())
		http.Error(w, "Could not write response", http.StatusInternalServerError)
		return
	}
}
