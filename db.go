package main

import (
	"database/sql"

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
			book_title VARCHAR(255),
			book_author VARCHAR(255),
			FOREIGN KEY (collection_id) REFERENCES collections(id),
			FOREIGN KEY (book_title, book_author) REFERENCES books(title, author)
		);
    `)

	if err != nil {
		return err
	}
	return nil
}
