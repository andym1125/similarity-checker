package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	buildHandler = http.FileServer(http.Dir("../frontend/build"))
	hashes       [][]byte
)

type CtphResponse struct {
	Percentage float64
	Substrings []Substring
	Control    []byte
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

func processHandler(w http.ResponseWriter, r *http.Request) {

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

	//Process request
	var processedRes CtphResponse

	hashOfTxt := hash([]byte(body.Text))
	allCommonSubstrs := getAllCommonSubstrings(hashOfTxt)

	processedRes.Substrings = allCommonSubstrs
	processedRes.Control = hashOfTxt
	processedRes.Percentage = float64(int(10000*getPercentage(hashOfTxt, allCommonSubstrs))) / 10000

	fmt.Println("Hash of input:", hashOfTxt)

	//Marshal response
	marshalByte, marshalErr := json.Marshal(processedRes)
	if marshalErr != nil {
		panic(marshalErr)
	}

	w.Write(marshalByte)
}

func getAllCommonSubstrings(control []byte) []Substring {

	var unparsedSubstrs []Substring
	for _, h := range hashes {
		unparsedSubstrs = append(unparsedSubstrs, getCommonSubstrings(control, h)...)
	}

	return prune(unparsedSubstrs)
}

/* ********** File Handling ********** */

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

	hashes = append(hashes, hash(data))
}
