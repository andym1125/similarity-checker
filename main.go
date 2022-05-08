package main

import (
	"fmt"
	_ "fmt"
	"log"
	"net/http"
)

var (
	buildHandler = http.FileServer(http.Dir("frontend/build"))
)

func main() {

	r := http.NewServeMux()
	r.HandleFunc("/", index)

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func index(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		buildHandler.ServeHTTP(w, r)
		//http.ServeFile(w, r, "frontend/build/index.html")
		fmt.Println("Get req", r.URL)
	} else if r.Method == "POST" {

		if err := r.ParseForm(); err != nil {
			fmt.Println("ParseForm() err: ", err)
			return
		}

		fmt.Println("Hash: ", hash([]uint8(r.FormValue("txt"))))
		http.ServeFile(w, r, "frontend/build/index.html")
	}
}
