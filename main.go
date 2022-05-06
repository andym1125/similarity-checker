package main

import (
	"fmt"
	_ "fmt"
	"log"
	"net/http"
)

func main() {

	r := http.NewServeMux()

	r.HandleFunc("/check", index)
	buildHandler := http.FileServer(http.Dir("frontend/build"))
	r.Handle("/", buildHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
	}

	fmt.Println("Listening on port 8080...")
	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/build/index.html")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

// func main() {
// 	handleRequests()
// }
