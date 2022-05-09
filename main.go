package main

import (
	"fmt"
	_ "fmt"
	"log"
	"net/http"
)

var (
	buildHandler = http.FileServer(http.Dir("frontend/build"))
	hashes       [][]byte
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

		//Allows for commands
		switch r.FormValue("txt") {
		case ".print":
			fmt.Println(hashes)
		case ".clear":
			hashes = [][]byte{}
			fmt.Println("Cleared storage.")
		default:
			h := hash([]uint8(r.FormValue("txt")))

			percent := checkPlagarism(h)

			hashes = append(hashes, h)
			fmt.Println("Hash:", h, percent)
		}

		http.ServeFile(w, r, "frontend/build/index.html")
	}
}

func checkPlagarism(h []byte) float64 {

	calcCommonSubstrings()

	substrs := []Substring{}
	for _, a := range hashes {
		substrs = append(substrs, getCommonSubstrings(h, a)...)
	}

	s := getPercentage(h, substrs)

	return s
}
