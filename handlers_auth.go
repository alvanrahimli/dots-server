package main

import (
	"encoding/json"
	"fmt"
	"github.com/alvanrahimli/dots-server/dataaccess"
	"github.com/alvanrahimli/dots-server/models"
	"github.com/alvanrahimli/dots-server/utils"
	"net/http"
	"time"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Printf("Request received for URL %s", r.URL)

	if err := r.ParseForm(); err != nil {
		_, err := fmt.Fprintf(w, "ParseForm() err: %v", err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ErrLogger.Println(err.Error())
			return
		}
	}

	username := r.FormValue("name")
	password := r.FormValue("password")
	email := r.FormValue("email")

	nameValidationResult, nameValidated := utils.ValidateName(username)
	passValidationResult, passValidated := utils.ValidatePassword(password)
	emailValidationResult, emailValidated := utils.ValidateEmail(email)

	response := models.NewHttpResponse()

	if !nameValidated || !passValidated || !emailValidated {
		w.WriteHeader(http.StatusBadRequest)
		response.Code = 1
		response.Message = "Validation error occurred"
		response.Data = map[string]string{
			"name":     nameValidationResult,
			"password": passValidationResult,
			"email":    emailValidationResult,
		}
	} else {
		// Register new user
		db := getDbInstance()
		//goland:noinspection GoUnhandledErrorResult
		defer db.Close()

		insertUserRaw := `INSERT INTO Users (Username, Email, Password) VALUES (?, ?, ?)`
		statement, stmtErr := db.Prepare(insertUserRaw)
		if stmtErr != nil {
			ErrLogger.Println(stmtErr.Error())
			http.Error(w, stmtErr.Error(), http.StatusInternalServerError)
			return
		}

		_, dbErr := statement.Exec(username, email, utils.HashPassword(password))
		if dbErr != nil {
			ErrLogger.Println(dbErr.Error())
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		response.Code = 0
		response.Message = fmt.Sprintf("User created with name '%s'", username)
		response.Data = map[string]string{
			"username": username,
			"email":    email,
		}
	}

	responseJson, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		ErrLogger.Println(jsonErr.Error())
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Code == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	if _, err := w.Write(responseJson); err != nil {
		ErrLogger.Println(err.Error())
		return
	}

	InfoLogger.Println("Request finished")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Printf("Request received for URL %s", r.URL)

	// Parse form
	if err := r.ParseForm(); err != nil {
		_, err := fmt.Fprintf(w, "ParseForm() err: %v", err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ErrLogger.Println(err.Error())
			return
		}
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	db := getDbInstance()
	//goland:noinspection GoUnhandledErrorResult
	defer db.Close()

	// Get user by email
	user, getUserErr := dataaccess.FindUserByEmail(db, email)
	if getUserErr != nil {
		ErrLogger.Println(getUserErr.Error())
		http.Error(w, fmt.Sprintf("Could not find user with email '%s'", email), http.StatusNotFound)
		return
	}

	// Check user's password
	if user.Password != utils.HashPassword(password) {
		WarnLogger.Printf("User '%s' entered invalid password", email)
		http.Error(w, fmt.Sprintf("User '%s' entered invalid password", email), http.StatusUnauthorized)
		return
	}

	// Delete existing tokens
	deleteOldTokensQuery := `DELETE FROM AuthTokens WHERE UserId = ?`
	statement, stmtErr := db.Prepare(deleteOldTokensQuery)
	if stmtErr != nil {
		ErrLogger.Println(stmtErr.Error())
		http.Error(w, stmtErr.Error(), http.StatusInternalServerError)
		return
	}

	_, execErr := statement.Exec(user.Id)
	if execErr != nil {
		ErrLogger.Println(execErr.Error())
		http.Error(w, "Could not delete existing tokens", http.StatusInternalServerError)
		return
	}

	// Insert new token
	token := models.AuthToken{
		Id:             0,
		UserId:         user.Id,
		Token:          utils.GenerateToken(),
		ExpirationDate: time.Now().Add(time.Hour * 24 * 31).Format(models.DatetimeLayout), // a month,
	}

	tokenErr := dataaccess.AddToken(&token, db)
	if tokenErr != nil {
		// TODO: token != nil
	}

	userDto := map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}

	response := models.HttpResponse{
		Code:    0,
		Message: "User authenticated successfully",
		Data: map[string]interface{}{
			"token":           token.Token,
			"expiration_date": token.ExpirationDate,
			"user":            userDto,
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
