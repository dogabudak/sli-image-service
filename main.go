package main

import (
	"io/ioutil"
	"net/http"
    "log"
    "math/rand"
    "time"
)

func main() {
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/photo", handler)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    files, err := ioutil.ReadDir("./images")
    if err != nil {
        log.Fatal(err)
    }
    rand.Seed(time.Now().Unix())
    n := rand.Int() % len(files)
    randomFileAddress := "./images/" + files[n].Name()
	fileBytes, err := ioutil.ReadFile(randomFileAddress)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}