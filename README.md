# Book Management REST API (Golang)

A simple RESTful API built using **Golang** and **Gorilla Mux** that performs
CRUD (Create, Read, Update, Delete) operations on books using an in-memory data store.

This project is intended for learning REST APIs, HTTP handling, and routing in Go.

---

## Features

- Create a new book
- Get all books
- Get a book by ID
- Update a book by ID
- Delete a book by ID
- JSON-based request and response handling

---

## Tech Stack

- **Language:** Go (Golang)
- **Router:** Gorilla Mux
- **Data Storage:** In-memory slice (no database)
- **Protocol:** HTTP / REST

---

## Project Structure

```text
restapi/
│
├── main.go          # Entry point of the application
├── handler.go       # Contains all API handlers
└── README.md        # Project documentation

Data Model
Book
{
  "id": 1,
  "title": "Go Programming",
  "author": {
    "auid": 101,
    "aname": "Alan A. A. Donovan"
  },
  "price": 499.99,
  "orders": {
    "oid": 5001,
    "customername": "Abhi",
    "totalprice": 999.98
  }
}

API Endpoints
Create Book

POST /books

Request Body
{
  "title": "Go in Action",
  "author": {
    "auid": 1,
    "aname": "William Kennedy"
  },
  "price": 399.99,
  "orders": {
    "oid": 101,
    "customername": "John",
    "totalprice": 399.99
  }
}

Response
{
  "book": { ... },
  "message": "Book Created Successfully"
}

Get All Books

GET /books

Response
[
  { ... },
  { ... }
]

Get Book by ID

GET /books/{id}

Example
GET /books/1

Update Book

PUT /books/{id}

Request Body
{
  "title": "Advanced Go",
  "author": {
    "auid": 2,
    "aname": "Rob Pike"
  },
  "price": 599.99,
  "orders": {
    "oid": 202,
    "customername": "Alice",
    "totalprice": 599.99
  }
}

Delete Book

DELETE /books/{id}

Response
{
  "message": "Deleted successfully"
}

▶Running the Application
1 Clone the Repository
git clone https://github.com/your-username/book-rest-api.git
cd book-rest-api

2️ Install Dependencies
go mod tidy

3️ Run the Server
go run main.go

4️ Server Starts At
http://localhost:8080

Testing the API

You can test the endpoints using:

cURL

Thunder Client (VS Code)

Example:

curl http://localhost:8080/books