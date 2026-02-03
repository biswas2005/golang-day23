package restapi

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Price  float64 `json:"price"`
	Author Author  `json:"author"`
	Orders Orders  `json:"orders"`
}
type Author struct {
	Auid  int    `json:"auid"`
	AName string `json:"aname"`
}
type Orders struct {
	Oid          sql.NullInt64   `json:"oid"`
	CustomerName sql.NullString  `json:"customername"`
	TotalPrice   sql.NullFloat64 `json:"totalprice"`
}

func postBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"INSERT INTO authors (aname) VALUES (?)",
		book.Author.AName,
	)
	if err != nil {
		http.Error(w, `{"error":"Author insert failed."}`, http.StatusInternalServerError)
		return
	}
	autherID, _ := res.LastInsertId()

	bookRes, err := db.Exec(
		"INSERT INTO books(title ,price, author_id)VALUES (?,?,?)",
		book.Title, book.Price, autherID,
	)
	if err != nil {
		http.Error(w, `{"error":"Book insert failed"}`, http.StatusInternalServerError)
		return
	}
	bookID, _ := bookRes.LastInsertId()

	_, err = db.Exec(
		"INSERT INTO orders (book_id, customer_name, total_price) VALUES (?, ?, ?)",
		bookID, book.Orders.CustomerName, book.Orders.TotalPrice,
	)
	if err != nil {
		http.Error(w, `{"error":"order insert failed"}`, http.StatusInternalServerError)
		return
	}
	book.ID = int(bookID)
	book.Author.Auid = int(autherID)

	w.WriteHeader(http.StatusCreated)

	err1 := json.NewEncoder(w).Encode(map[string]interface{}{
		"book":    book,
		"message": "Book Created Successfully",
	})
	if err1 != nil {
		http.Error(w, `{"error":"Failed to Encode Response"}`, http.StatusInternalServerError)
	}

}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`
		SELECT b.id, b.title, b.price,
		       a.auid, a.aname,
		       o.oid, o.customer_name, o.total_price
		FROM books b
		JOIN authors a ON b.author_id = a.auid
		LEFT JOIN orders o ON o.book_id = b.id
	`)
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Price,
			&book.Author.Auid,
			&book.Author.AName,
			&book.Orders.Oid,
			&book.Orders.CustomerName,
			&book.Orders.TotalPrice,
		); err != nil {
			http.Error(w, "Row scan failed:"+err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "JSON encode failed:"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID Format", http.StatusBadRequest)
		return
	}

	var book Book
	err = db.QueryRow(`
	SELECT b.id,b.title,b.price,
		a.auid,a.aname,
		o.oid,o.customer_name,o.total_price
		FROM books b
		JOIN authors a ON b.author_id=a.auid 
		LEFT JOIN orders o ON o.book_id=b.id
		WHERE b.id=?
	`, id).Scan(
		&book.ID,
		&book.Title,
		&book.Price,
		&book.Author.Auid,
		&book.Author.AName,
		&book.Orders.Oid,
		&book.Orders.CustomerName,
		&book.Orders.TotalPrice,
	)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Book not Found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"Query Failed"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "JSON encode failed", http.StatusInternalServerError)
		return
	}
}

func putBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id, err1 := strconv.Atoi(params["id"])
	if err1 != nil {
		http.Error(w, "Invalid ID Format", http.StatusBadRequest)
		return
	}
	var updated Book

	err := json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	_, err0 := db.Exec(`
	UPDATE books SET title=?,price=? WHERE id=?`, updated.Title, updated.Price, id)

	if err0 != nil {
		http.Error(w, `{"error":"Update Failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Updated Successfully.",
	})
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID Format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM orders WHERE book_id=?", id)
	if err != nil {
		http.Error(w, "Delete failed.", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("DELETE FROM books WHERE id=?", id)
	if err != nil {
		http.Error(w, "Delete failed.", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Deleted Successfully.",
	})
}

func Handler() {

	router := mux.NewRouter()

	router.HandleFunc("/books", postBook).Methods("POST")
	router.HandleFunc("/books", getAllBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", putBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
