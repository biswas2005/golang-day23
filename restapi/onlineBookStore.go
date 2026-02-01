package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book represents the data model for our API.
// It contains nested Author and Orders structs to demonstrate relationships.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author struct {
		Auid  int    `json:"auid"`
		AName string `json:"aname"`
	} `json:"author"`
	Price  float64 `json:"price"`
	Orders struct {
		Oid          int     `json:"oid"`
		CustomerName string  `json:"customername"`
		TotalPrice   float64 `json:"totalprice"`
	} `json:"orders"`
}

// books acts as our in-memory database.
var books []Book

// postBook() handles POST /books
// It creates a new book entry and assigns a unique ID.
func postBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	//Decode JSON request body into Book struct.
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	//Generate new ID by finding max existing ID
	maxID := 0
	for _, bs := range books {
		if bs.ID > maxID {
			maxID = bs.ID
		}
	}
	book.ID = maxID + 1
	// Append new book to slice
	books = append(books, book)
	w.WriteHeader(http.StatusCreated)
	// Respond with created book and success message
	err1 := json.NewEncoder(w).Encode(map[string]interface{}{
		"book":    book,
		"message": "Book Created Successfully",
	})
	if err1 != nil {
		http.Error(w, `{"error":"Failed to Encode Response"}`, http.StatusInternalServerError)
	}

}

// getBooks() handles GET /books
// It returns all books in the slice.
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// getBook() handles GET /books/{id}
// It returns a single book by ID.
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Extract ID from URL params
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid ID Format",
		})
		return
	}
	// Search for book by ID
	for _, book := range books {
		if book.ID == id {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(book)
			if err != nil {
				http.Error(w, `{"error":"Failed to Encode"}`, http.StatusInternalServerError)
			}
			return
		}
	}
	//If not found
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("Could not find ID")
}

// putBook() handles PUT /books/{id}
// Updates an existing book
func putBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//Extract id from URL params
	idx, err1 := strconv.Atoi(params["id"])
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid ID Format",
		})
		return
	}
	var updated Book
	//Decode updated book data from Request body
	err := json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	exists := false
	//find and update book in Slice
	for id, book := range books {
		if book.ID == idx {
			updated.ID = idx
			books[id] = updated
			exists = true
			break
		}
	}
	//if Book id not found
	if !exists {
		http.Error(w, `{"error":"Could not find ID"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//send Success response
	errr := json.NewEncoder(w).Encode(map[string]interface{}{
		"book":    updated,
		"message": "Updated Successfully",
	})
	if errr != nil {
		http.Error(w, `{"error":"Failed to Encode"}`, http.StatusInternalServerError)
	}
}

// deleteBook() handles DELETE /books/{id}
// Removes a book by ID
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//extract id from URL params
	idx, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid ID Format",
		})
		return
	}
	exists := false
	//find and delete book from slice
	for i, val := range books {
		if val.ID == idx {
			//remove element at index i
			books = append(books[:i], books[i+1:]...)
			exists = true
			break
		}
	}
	//If id not Found
	if !exists {
		http.Error(w, `{"error":"ID does not exist"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//send success response
	err1 := json.NewEncoder(w).Encode(map[string]string{
		"message": "Deleted successfully",
	})
	if err1 != nil {
		http.Error(w, `{"error":"Failed to Encode"}`, http.StatusInternalServerError)
	}
}

// Handler() sets up routes and starts the HTTP server
func Handler() {
	//Create a nes Gorilla Mux router
	router := mux.NewRouter()
	//route definations
	router.HandleFunc("/books", postBook).Methods("POST")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", putBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	//Start HTTP server on port 8080
	http.ListenAndServe(":8080", router)
}
