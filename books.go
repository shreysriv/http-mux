package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)
var mu sync.Mutex
func sequence(initValue int) func() int {
	i := initValue

	return func() int {
		i++
		return i
	}
}

// Book represents a single book
type Book struct {
	ID          int
	Title       string
	Author      string
	ISBN        string
	Description string
	Price       float64
}

var nextID = sequence(0)

var books = map[int]Book{
	1: Book{ID: nextID(), Title: "The C Book", Author: "Dennis Ritchie"},
	2: Book{ID: nextID(), Title: "C++", Author: "Bjarne Stroustrop"},
}

// GET /books/{id}
func bookShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, _ := strconv.Atoi(vars["id"])
	log.Println("Getting bookID: ", bookID)
	book := books[bookID]

	json.NewEncoder(w).Encode(book)
}

// GET /books
func booksIndexHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

// POST /books
func booksCreateHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}
	mu.Lock()
	books[len(books)+1] = book
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(book)
}

// DELETE /books
func booksDestroyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, _ := strconv.Atoi(vars["id"])
	log.Println("Deleting BookId: ", bookID)
	delete(books, bookID)
}

// PUT /books
func booksUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	vars := mux.Vars(r)
	bookID, _ := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}
	mu.Lock()
	log.Println("Updating BookId: ", bookID)
	books[bookID] = book
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(book)
}

// PATCH /books
func booksPatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, _ := strconv.Atoi(vars["id"])
	book := books[bookID]
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}
	mu.Lock()
	log.Println("Updating BookId: ", bookID)
	books[bookID] = book
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(book)
}

func main() {
	port := "9000"
	r := mux.NewRouter()

	r.HandleFunc("/books", booksIndexHandler).Methods("GET")
	r.HandleFunc("/books/{id}", bookShowHandler).Methods("GET")
	r.HandleFunc("/books", booksCreateHandler).Methods("POST")
	r.HandleFunc("/books/{id}", booksUpdateHandler).Methods("PUT")
	r.HandleFunc("/books/{id}", booksPatchHandler).Methods("PATCH")
	r.HandleFunc("/books/{id}", booksDestroyHandler).Methods("DELETE")
	log.Println("Server is running on port: ", port)

	http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r))
}
