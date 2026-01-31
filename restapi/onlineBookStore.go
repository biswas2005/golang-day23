package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

var books []Book

func postBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	maxID := 0
	for _, bs := range books {
		if bs.ID > maxID {
			maxID = bs.ID
		}
	}
	book.ID = maxID + 1
	books = append(books, book)
	w.WriteHeader(http.StatusCreated)

	err1 := json.NewEncoder(w).Encode(map[string]interface{}{
		"book":    book,
		"message": "Book Created Successfully",
	})
	if err1 != nil {
		http.Error(w, `{"error":"Failed to Encode Response"}`, http.StatusInternalServerError)
	}

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid ID Format",
		})
		return
	}
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
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("Could not find ID")
}

func putBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idx, err1 := strconv.Atoi(params["id"])
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid ID Format",
		})
		return
	}
	var updated Book

	err := json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	exists := false
	for id, book := range books {
		if book.ID == idx {
			updated.ID = idx
			books[id] = updated
			exists = true
			break
		}
	}
	if !exists {
		http.Error(w, `{"error":"Could not find ID"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	errr := json.NewEncoder(w).Encode(map[string]interface{}{
		"book":    updated,
		"message": "Updated Successfully",
	})
	if errr != nil {
		http.Error(w, `{"error":"Failed to Encode"}`, http.StatusInternalServerError)
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idx, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid ID Format",
		})
		return
	}
	exists := false
	for i, val := range books {
		if val.ID == idx {
			books = append(books[:i], books[i+1:]...)
			exists = true
			break
		}
	}
	if !exists {
		http.Error(w, `{"error":"ID does not exist"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err1 := json.NewEncoder(w).Encode(map[string]string{
		"message": "Deleted successfully",
	})
	if err1 != nil {
		http.Error(w, `{"error":"Failed to Encode"}`, http.StatusInternalServerError)
	}
}

func Handler() {

	router := mux.NewRouter()

	router.HandleFunc("/books", postBook).Methods("POST")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", putBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
