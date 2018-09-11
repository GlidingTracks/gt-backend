package main

import (
	"firebase.google.com/go"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

func main() {
	initializeFirebase()

	r := mux.NewRouter()

	r.HandleFunc("/", startPage)

	logrus.Fatal(http.ListenAndServe(":8080", r))
}

//Redirect here is url: localhost:8080 is supplied
func startPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from go - Gliding tracks\n")
}

func initializeFirebase() {
	opt := option.WithCredentialsFile("gt-backend-8b9c2-firebase-adminsdk-0t965-80679b9b72.json")
	/*app*/ _, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
}
