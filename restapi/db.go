package restapi

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDB() {
	var err error

	dsn := "root:root@tcp(localhost:3306)/bookstore?parseTime=true"

	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB Open Error:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("DB Connection Error:", err)
	}
	log.Println("MySQL Connected Successfully.")
}
