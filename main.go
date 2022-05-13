package main

import (
	"fmt"
	_ "fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	buildHandler = http.FileServer(http.Dir("frontend/build"))
	hashes       [][]byte
	doStoreHash  = false
)

func main() {

	r := http.NewServeMux()
	r.HandleFunc("/", index)

	loadReferences()

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

		case ".store":
			doStoreHash = !doStoreHash
			fmt.Println("Store Hashes:", doStoreHash)

		default:
			h := hash([]uint8(r.FormValue("txt")))
			if doStoreHash {
				directStoreHash(h)
				fmt.Println("Storing hash...")
			}

			percent := checkPlagarism(h)
			fmt.Println("Hash:", h, percent)
		}

		http.ServeFile(w, r, "frontend/build/index.html")
	}
}

func checkPlagarism(h []byte) float64 {

	substrs := []Substring{}
	for _, a := range hashes {
		substrs = append(substrs, getCommonSubstrings(h, a)...)
	}

	return getPercentage(h, substrs)
}

func storeHash(f []byte) []byte {

	h := hash(f)
	hashes = append(hashes, h)
	return h
}

func directStoreHash(h []byte) {
	hashes = append(hashes, h)
}

func loadReferences() {

	loadFile("files/rollinghash.txt")
	loadFile("files/shakespeare.txt")
}

func loadFile(path string) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	storeHash(data)
}
