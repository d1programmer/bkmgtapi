package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	Edition       string `json:"edition"`
	Description   string `json:"description"`
	Genre         string `json:"genre"`
}

type Collection struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Configure MySQL database connection
	var err error
	db, err = connectDB()
	if err != nil {
		fmt.Print("Error connecting to SQL Database\n\n\n")
		log.Fatal("Error connecting to the database:", err)
	}

	fmt.Print("Successfully connected to MySQL DB\n\n\n")
	defer db.Close()

	err = initializeDatabaseTables(db)
	if err != nil {
		fmt.Print("Error initializing database tables\n\n\n")
		log.Fatal("Error initializing database tables:", err)
	}

	fmt.Println("Book Management API is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("Error starting server\n\n\n")
		log.Fatal("Error starting server:", err)
	}
}
