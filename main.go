package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

type CollectionEntry struct {
	Collection_id int
	Book_id       int
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

	router.HandleFunc("/", welcomeHandler)
	router.HandleFunc("/books", booksHandler)                                                               // get and post
	router.HandleFunc("/books/{bookid:[0-9]+}", booksPutAndDeleteHandler)                                   // edit a book
	router.HandleFunc("/collections", collectionsHandler)                                                   // get and post (list collections, create collections)
	router.HandleFunc("/collections/{collectionid:[0-9]+}", collectionsGetPostAndDeleteHandler)             // get specific collection, add to collection, delete entire collection
	router.HandleFunc("/collections/{collectionid:[0-9]+}/books/{bookid:[0-9]+}", collectionsDeleteHandler) // delete book from collection

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
	method := r.Method

	switch method {
	case http.MethodGet:
		getBooksHandler(w, r)
	case http.MethodPost:
		postBooksHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, POST, DELETE")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func booksPutAndDeleteHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case http.MethodPut:
		putBooksHandler(w, r) // NOT YET IMPLEMENTED
	case http.MethodDelete:
		deleteBooksHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, POST, DELETE")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func collectionsHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case http.MethodGet:
		getCollectionsHandler(w, r)
	case http.MethodPost:
		postCollectionsHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func collectionsGetPostAndDeleteHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case http.MethodGet:
		getSpecificCollectionsHandler(w, r) // NOT YET IMPLEMENTED
	case http.MethodPost:
		addToCollectionHandler(w, r)
	case http.MethodDelete:
		deleteCollectionHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func collectionsDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := getBooks()

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Encode books into JSON
	jsonResponse, err := json.Marshal(books)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(books) == 0 {
		jsonResponse = []byte("[]")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func postBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)

	if err != nil {
		http.Error(w, "Error decoding JSON request", http.StatusBadRequest)
		return
	}

	// Insert book into the database
	bookCreated, err := insertBookIntoDB(newBook)
	if err != nil {
		http.Error(w, "Error inserting book into database", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	bookJson, err := json.Marshal(bookCreated)

	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bookJson)
	fmt.Fprintln(w, "\nBook added successfully!")
}

func putBooksHandler(w http.ResponseWriter, r *http.Request) {
	// TO BE IMPLEMENTED
}

func deleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	bookIDStr := vars["bookid"]

	bookID, err := strconv.Atoi(bookIDStr)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = deleteBook(bookID)

	if err.Error() == struct{ error }{fmt.Errorf("404")}.Error() {
		http.Error(w, "Not found", http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "Book deleted successfully!")
}

func getCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collections, err := getCollections()

	if err != nil {

	}
	// Encode collections into JSON
	jsonResponse, err := json.Marshal(collections)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if len(collections) == 0 {
		jsonResponse = []byte("[]")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func postCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newCollection Collection
	err := json.NewDecoder(r.Body).Decode(&newCollection)

	if err != nil {
		http.Error(w, "Error decoding JSON request", http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = createCollection(newCollection)

	if err != nil {
		http.Error(w, "Error inserting collection into database", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	collectionCreated, err := json.Marshal(newCollection)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(collectionCreated)
	fmt.Fprintln(w, "\nCollection added successfully!")
}

func getSpecificCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	// TO BE IMPLEMENTED
}

func addToCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	collectionIDStr := vars["collectionid"]

	collectionID, err := strconv.Atoi(collectionIDStr)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var requestData struct {
		BookID int `json:"bookid"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)

	if err != nil {
		http.Error(w, "Error decoding JSON request", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = addBookToCollection(collectionID, requestData.BookID)

	var c CollectionEntry
	c.Collection_id = collectionID
	c.Book_id = requestData.BookID

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	collectionCreated, err := json.Marshal(c)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(collectionCreated)

	fmt.Fprintln(w, "\nBook added to collection successfully!")
}

func deleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	collectionIDStr := vars["collectionid"]

	collectionID, err := strconv.Atoi(collectionIDStr)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = deleteCollection(collectionID)

	if err.Error() == struct{ error }{fmt.Errorf("404")}.Error() {
		http.Error(w, "Not found", http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "Collection deleted successfully!")
}
