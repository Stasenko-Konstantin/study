package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func handler(h string, w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("assets/" + h + ".html")
	if err != nil {
		fmt.Fprintf(w, "Server error: "+err.Error())
	}
	fmt.Fprintf(w, string(file))
}

func main() {
	fmt.Println("server start")
	r := mux.NewRouter()
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) { handler("index", w, r) }).Methods("GET")
	r.HandleFunc("/games",
		func(w http.ResponseWriter, r *http.Request) { handler("games", w, r) }).Methods("GET")
	r.HandleFunc("/article",
		func(w http.ResponseWriter, r *http.Request) { handler("article", w, r) }).Methods("GET")
	staticFileDir := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDir))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	http.ListenAndServe(":8080", r)
}
