package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	fmt.Fprintf(w, "POST request successful! \n")
	book := r.FormValue("book")
	author := r.FormValue("author")

	// w is the IO writer to send sponse to our website
	fmt.Fprintf(w, "Name = %s\n", book)
	fmt.Fprintf(w, "Author = %s\n", author)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

func main() {
	// GO will auto look at file index in the ./static to create server
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting server at port 8000\n")
	// ListenAndServe is the way to start server
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
