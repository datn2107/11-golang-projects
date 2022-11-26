package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	// `json:"id"` is a struct tag
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var books []Book
var n_books int

func deleteBookInternal(id string) {
	for idx, item := range books {
		if item.ID == id {
			books = append(books[:idx], books[idx+1:]...)
			break
		}
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	deleteBookInternal(params["id"])

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	n_books++
	book.ID = strconv.Itoa(n_books)
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	deleteBookInternal(params["id"])
	book.ID = params["id"]
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func main() {
	books = append(books, Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &Author{"Dan", "Brown"}})
	books = append(books, Book{ID: "2", Isbn: "438241", Title: "Book Two", Author: &Author{"Dan", "White"}})
	n_books = 2

	r := mux.NewRouter()

	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Printf("Start server at port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}
