package main

import (
	"net/http"
	"os"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Printf("PING! received")

	if r.URL.Path != "/ping" {
		ErrLogger.Println("Url does not match: /ping")
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("DOTS-CLI-VERSION") == "" {
		ErrLogger.Printf("Header DOTS-CLI-VERSION is not defined (remote: %s)", r.RemoteAddr)
		http.Error(w, "", http.StatusTeapot)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Server", os.Getenv("REGISTRY_DOMAIN"))
	_, err := w.Write([]byte("pong"))
	if err != nil {
		ErrLogger.Println(err.Error())
	}

	InfoLogger.Println("Request finished")
}
