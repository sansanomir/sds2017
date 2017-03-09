package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	// Redirect the incoming HTTP request. Note that "127.0.0.1:8081" will only work if you are accessing the server from your local machine.
	http.Redirect(w, r, "https://127.0.0.1:8081"+r.RequestURI, http.StatusMovedPermanently)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/read" {
		fmt.Fprintf(w, r.RequestURI)
		file, err := os.Open("datos.txt")
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()
		str, err := ioutil.ReadAll(file)
		fmt.Println(string(str))
	} else if r.RequestURI == "/write" {
		fmt.Fprintf(w, r.RequestURI)
		fileW, err := os.Create("result.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer fileW.Close()
		fmt.Fprintf(fileW, "user: Pepe\nmasterKey: asdf\n")
	} else {
		fmt.Fprintf(w, "I don't know")
	}
}

func main() {
	http.HandleFunc("/", handler)
	// Start the HTTPS server in a goroutine
	go http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
	// Start the HTTP server and redirect all incoming connections to HTTPS
	http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))
}
