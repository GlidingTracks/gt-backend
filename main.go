package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", startPage)
	r.HandleFunc("/upload", uploadFilePage).Methods("POST")

	logrus.Fatal(http.ListenAndServe(":8080", r))
}

//Redirect here is url: localhost:8080 is supplied
func startPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from go - Gliding tracks\n")
}

// Upload and save a file to the filesystem
func uploadFilePage(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	if uid == "" {
		logrus.Error("No uid supplied in request")
		http.Error(w, "No uid supplied in request", http.StatusBadRequest)
		return
	}

	httpCode, err := ProcessUploadRequest(r, uid)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), httpCode)
	}
}
