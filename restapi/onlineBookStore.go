package restapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
var count Book

func fileHandle() (*os.File, error) {
	file, err := os.OpenFile("BookStore.file", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func loadStore() {
	books = nil
	file, err := os.Open("BookStore.file")
	if err != nil {
		fmt.Println("Error Opening File", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var b Book
		err := json.Unmarshal([]byte(line), &b)
		if err != nil {
			fmt.Println("Error Unmarshalling", err)
			return
		}
		books = append(books, b)
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	loadStore()
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	exists := false
	for _, bs := range books {
		if book.ID == bs.ID {
			count.ID++
			exists = true
			break
		}
	}
	if !exists {
		book.ID = len(books) + 1
		books = append(books, book)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	}
}
func Handler() {
	fileHandle()
	router := mux.NewRouter()

	router.HandleFunc("/books", addBook).Methods("POST")

	http.ListenAndServe(":8080", router)
}
