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
