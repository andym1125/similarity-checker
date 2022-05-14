package main

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	buildHandler = http.FileServer(http.Dir("../frontend/build"))
	hashes       [][]byte
	doStoreHash  = false
)

type CtphResponse struct {
	Percentage float64
	Substrings []Substring
	Control    string
}

type TestResponse struct {
	Test string `json:"test"`
}

type CtphRequest struct {
	Text string `json:"text"`
}

func main() {

	r := http.NewServeMux()
	r.HandleFunc("/", buildHandler.ServeHTTP)
	r.HandleFunc("/process", processHandler)

	loadReferences()

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func index(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		buildHandler.ServeHTTP(w, r)
		http.ServeFile(w, r, "frontend/build/index.html")
		fmt.Println("Get req", r.URL)
	} else if r.Method == "POST" {

		// if err := r.ParseForm(); err != nil {
		// 	fmt.Println("ParseForm() err: ", err)
		// 	return
		// }

		// //Allows for commands
		// switch r.FormValue("txt") {
		// case ".print":
		// 	fmt.Println(hashes)

		// case ".clear":
		// 	hashes = [][]byte{}
		// 	fmt.Println("Cleared storage.")

		// case ".store":
		// 	doStoreHash = !doStoreHash
		// 	fmt.Println("Store Hashes:", doStoreHash)

		// default:
		// 	h := hash([]uint8(r.FormValue("txt")))
		// 	if doStoreHash {
		// 		directStoreHash(h)
		// 		fmt.Println("Storing hash...")
		// 	}

		// 	percent := checkPlagarism(h)
		// 	fmt.Println("Hash:", h, percent)
		// }

		http.ServeFile(w, r, "../frontend/build/index.html")
	}
}

func processHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Anything")

	//Enable cors
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//Read the body
	bodyByte, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}

	//Unmarshal body
	var body CtphRequest
	unmarshalErr := json.Unmarshal(bodyByte, &body)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}

	fmt.Println("Request received: ", body.Text)

	//Marshal response
	processedRes := TestResponse{"hello world"}
	marshalByte, marshalErr := json.Marshal(processedRes)
	if marshalErr != nil {
		panic(marshalErr)
	}

	w.Write(marshalByte)
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

	loadFile("../files/rollinghash.txt")
	loadFile("../files/shakespeare.txt")
}

func loadFile(path string) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	storeHash(data)
}
