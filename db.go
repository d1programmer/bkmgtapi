package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func connectDB() (*sql.DB, error) {
	// Replace with your MySQL database credentials and connection details
	var err error
	connectionString := "desmondtalton:pixelate09@tcp(localhost:3306)/book_management"

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initializeDatabaseTables(db *sql.DB) error {
	_, err := db.Exec("USE book_management;")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS collection_books;")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS books;")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS collections;")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS books (
			id INT AUTO_INCREMENT PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            author VARCHAR(255) NOT NULL,
            published_date DATE,
            edition VARCHAR(50),
            description TEXT,
            genre VARCHAR(100),
            UNIQUE (title, author)
        );
    `)
	if err != nil {
		return err
	}

	// Create the collections table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS collections (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE
        );
    `)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS collection_books (
			collection_id INT,
			book_id INT,
			PRIMARY KEY (collection_id, book_id),
			FOREIGN KEY (collection_id) REFERENCES collections(id),
			FOREIGN KEY (book_id) REFERENCES books(id)
		);
    `)

	if err != nil {
		return err
	}
	return nil
}

func getBooks() ([]Book, error) {

	rows, err := db.Query("SELECT * FROM books;")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.Edition, &book.Description, &book.Genre)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func insertBookIntoDB(book Book) (Book, error) {
	result, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		&book.Title, &book.Author, &book.PublishedDate, &book.Edition, &book.Description, &book.Genre)

	if err != nil {
		return Book{}, err
	}

	insertedID, err := result.LastInsertId()

	if err != nil {
		return Book{}, err
	}

	book.ID = int(insertedID)

	return book, nil
}

func getCollections() ([]Collection, error) {
	rows, err := db.Query("SELECT * FROM collections;")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var collection Collection
		err := rows.Scan(&collection.ID, &collection.Name)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

func createCollection(collection Collection) error {
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", collection.Name)
	return err
}

func deleteBook(bookID int) error {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id = ?)", bookID).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return struct{ error }{fmt.Errorf("404")}
	}

	_, err = db.Exec("DELETE FROM collection_books WHERE book_id = ?", bookID)

	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM books WHERE id = ?", bookID)

	if err != nil {
		return err
	}
	return nil
}

func addBookToCollection(collectionID, bookID int) error {
	_, err := db.Exec("INSERT INTO collection_books (collection_id, book_id) VALUES (?, ?)",
		collectionID, bookID)
	if err != nil {
		return err
	}
	return nil
}

func deleteCollection(collectionID int) error {

	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM collections WHERE id = ?)", collectionID).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return struct{ error }{fmt.Errorf("404")}
	}
	// Delete all entries in collection_books that reference the deleted collection

	_, err = db.Exec("DELETE FROM collection_books WHERE collection_id = ?", collectionID)

	if err != nil {
		return err
	}

	// Delete the collection from the collections table
	_, err = db.Exec("DELETE FROM collections WHERE id = ?", collectionID)

	if err != nil {
		return err
	}

	return nil
}
