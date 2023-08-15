package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	router.HandleFunc("/books", booksHandler)                                                   // get and post
	router.HandleFunc("/books/{bookid:[0-9]+}", booksPutAndDeleteHandler)                       // edit a book
	router.HandleFunc("/collections", collectionsHandler)                                       // get and post (list collections, create collections)
	router.HandleFunc("/collections/{collectionid:[0-9]+}", collectionsGetPostAndDeleteHandler) //get specific collection, add to collection, delete entire collection
	router.HandleFunc("/collections/{collectionid:[0-9]+}", collectionsDeleteHandler)           // delete book from collection

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server
	fmt.Println("Book Management API is running on port 8080...")
	log.Fatal(server.ListenAndServe())
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Book Management API!")
}

func booksHandler(w http.ResponseWriter, r *http.Request) {

}

func booksPutAndDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func collectionsHandler(w http.ResponseWriter, r *http.Request) {

}

func collectionsGetPostAndDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func collectionsDeleteHandler(w http.ResponseWriter, r *http.Request) {

}
