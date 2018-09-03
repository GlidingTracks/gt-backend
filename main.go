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

	logrus.Fatal(http.ListenAndServe(":8080", r))
}

//Redirect here is url: localhost:8080 is supplied
func startPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from go - Gliding tracks\n")
}
