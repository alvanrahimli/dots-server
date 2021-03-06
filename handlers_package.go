package main

import (
	"encoding/json"
	"fmt"
	"github.com/alvanrahimli/dots-server/dataaccess"
	"github.com/alvanrahimli/dots-server/models"
	"github.com/alvanrahimli/dots-server/utils"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
)

func addPackageHandler(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Printf("%s: Request received for URL %s", r.Method, r.URL)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		_, err := fmt.Fprintf(w, "ParseForm() err: %v", err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ErrLogger.Println(err.Error())
			return
		}
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
		WarnLogger.Printf("User '%s' token validation failed", userId)
		http.Error(w, "", http.StatusForbidden)
		return
	}

	user, userErr := dataaccess.FindUserById(db, userId)
	if userErr != nil {
		ErrLogger.Println(userErr.Error())
		http.Error(w, "Could not get user", http.StatusNotFound)
		return
	}

	archiveName := path.Join(models.ArchivesFolderRoot, user.Username, path.Base(fileHeader.Filename))
	// Check user's archive folder
	_, dirStatErr := os.Stat(path.Dir(archiveName))
	if dirStatErr != nil {
		if os.IsNotExist(dirStatErr) {
			err := os.MkdirAll(path.Dir(archiveName), os.ModePerm)
			if err != nil {
				ErrLogger.Println(err.Error())
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				return
			}
		} else {
			ErrLogger.Println(dirStatErr.Error())
			http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
			return
		}
	}

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

	// pack.Name = my-pack@first-user.registry.com
	pack := models.Package{
		Id:          0,
		Name:        fmt.Sprintf("%s@%s.%s", packageName, user.Username, models.RegistryDomain),
		Version:     packageVersion,
		ArchiveName: archiveName,
		UserId:      userId,
	}
	InfoLogger.Println(pack.Name)

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
			http.Error(w, "Could not add package", http.StatusBadRequest)
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

func getPackagesHandler(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Printf("Request received for URL %s", r.URL)

	packageName := mux.Vars(r)["name"]
	if packageName == "" {
		http.Error(w, "Package name not provided", http.StatusBadRequest)
		return
	}

	db := getDbInstance()
	defer db.Close()

	packages, getErr := dataaccess.GetPackages(packageName, db)
	if getErr != nil {
		ErrLogger.Println(getErr.Error())
		http.Error(w, "Error occurred while getting packages", http.StatusInternalServerError)
		return
	}

	if len(packages) == 0 {
		http.Error(w, "Could not find any package", http.StatusNotFound)
		return
	}

	response := models.HttpResponse{
		Code:    0,
		Message: fmt.Sprintf("Found %d packages", len(packages)),
		Data: map[string]interface{}{
			"Packages": packages,
		},
	}

	responseJson, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		ErrLogger.Println(jsonErr.Error())
		http.Error(w, "Could not marshal response dto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, writeErr := w.Write(responseJson)
	if writeErr != nil {
		ErrLogger.Println(writeErr.Error())
		http.Error(w, "Could not write response", http.StatusInternalServerError)
		return
	}
}
