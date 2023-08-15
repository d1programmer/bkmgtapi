# STEP ONE: DESCRIBE THE EXPECTED USER EXPERIENCE

The following is an expected user experience for the book management software, including the CLI Commands and their options.
The Command Line Interface will allow for users to be able to utilize create, read, update, and delete (CRUD) operations to manage a simple book management software.

The software will allow users to pass in command line prompts to update and modify the Book Management Software. Following the execution of a prompt, the API will respond with a status code and a description of what went right* [ or wrong (>.<) ]*
The API operations and details regarding them will be listed in the user-manual.txt file which is located in this repository. This manual will provide commands, optional flags, and examples of how to properly use them along with a response.

Users will be able to:
    -add a book to the system
    -list all books in the system
    -search for a specific book in the system
    -modify information about a book that is already in the system
    -create a collection
    -add a book to a collection
    -delete a collection
    -remove a book from a collection

# STEP TWO: DESCRIBE THE EXPECTED REST API

Here is an outline of the REST API for the Book Management Software, including supported methods, functionality, I/O data examples, and query parameters. 

BOOK MANAGEMENT SOFTWARE API DOCUMENTATION

BASE URL:
This software only runs locally on port8080.
http://localhost:8080

AUTHENTICATION:
This software uses MySQL to locally run a database.
user-manual.txt will provide details on how properly initialize a mySQL database on your local device for testing purposes. 

ENDPOINTS:

1. Book Management

    - Add a book

        URL: `/books`

        Method: POST
        Request Body: 
        {
            "title": "The Great Gatsby",
            "author": "F. Scott Fitzgerald",
            "published_date": "1925-04-10",
            "edition": "1st",
            "description": "A classic novel about the Jazz Age in America.",
            "genre": "Fiction"
        }
        Response:
        HTTP Status: 201 Created
        {
            "id": "1"
            "title": "The Great Gatsby",
            "author": "F. Scott Fitzgerald",
            "published_date": "1925-04-10",
            "edition": "1st",
            "description": "A classic novel about the Jazz Age in America.",
            "genre": "Fiction"
        }

    - List books

        URL: `/books`

        Method: GET
        Response: 
        HTTP Status: 200 OK
        [
            {
                "id": "1"
                "title": "The Great Gatsby",
                "author": "F. Scott Fitzgerald",
                "published_date": "1925-04-10",
                "edition": "1st",
                "description": "A classic novel about the Jazz Age in America.",
                "genre": "Fiction"
            },
            {
                "id": "2"
                "title": "The Alchemist",
                "author": "Paulo Coehlo",
                "published_date": "1988-07-24",
                "edition": "25th",
                "description": "An enchanting story about an Andalusian shepherd boy.",
                "genre": "Fiction"
            },
            ...
        ]

        Notes: 
            - The API will return a list of book objects in JSON format.
            - Each book object contains an `id`, `title`, `author`, `published_date`, `edition`, `description`, and `genre`.

        Error Handling:
            - 404 Not Found: If there are no books in the system.

    - Edit a book

        URL: `books/:book_id`

        Method: PUT
        Request Body:
        {
            "title": "The Great Gatsby",
            "author": "F. Scott Fitzgerald",
            "published_date": "1925-04-10",
            "edition": "1st",
            "description": "A description much different from the last.",
            "genre": "Fiction"
        }
        Response:
        HTTP Status: 200 OK
        {
            "id": 1
            "title": "The Great Gatsby",
            "author": "F. Scott Fitzgerald",
            "published_date": "1925-04-10",
            "edition": "1st",
            "description": "A description much different from the last.",
            "genre": "Fiction"
        }


        Notes:
            - Provide the book_id in the URL to specify the book you want to update.
            - The request body should contain the updated information for the book.
            - The API will return the updated book object in JSON format.
        
        Error Handling:
            - 400 Bad Request: If the request body is missing or holds improper information.
            - 404 Not Found: If the book is not found in the database.

    - Delete a book
        
        URL: `books/:book_id`

        Method: DELETE
        Response: 
        HTTP Status: 204 No Content

        Notes:
            - Provide the book_id in the URL to specify the book you want to delete.
        
        Error Handling:
            - 404 Not Found: If the specified book is not found in the database.

2. Collection Management

    - Create a collection

        URL: `/collections`

        Method: POST
        Request Body: 
        {
            "name": "Fiction"
        }
        Response:
        HTTP Status: 201 Created
        {
            "name": "Fiction"
        }

    - List collections
    
        URL: `/collections`

        Method: GET
        Response:
        HTTP Status: 200 OK
        [
            {
                "id": 1,
                "name": "Fiction"
            },
            {
                "id": 2,
                "name": "Non-fiction"
            },
            ...
        ]

        Notes: 
            - The API will return a list of collection objects in JSON format.
            - Each colleciton object contains an `id`, and `name`.

        Error Handling:
            - 404 Not Found: If there are no books in the system.

    - List books in a collection

        URL: `/collections/:collection_id`

        Method: GET
        Response:
        HTTP Status: 200 OK
        {
            "id": "1",
            "name": "Fiction",
            "books": [
                {
                    "id": 1,
                    "title": "The Great Gatsby",
                    "author": "F. Scott Fitzgerald",
                    "published_date": "1925-04-10",
                    "edition": "1st",
                    "description": "A classic novel about the Jazz Age in America.",
                    "genre": "Fiction"
                },
                {
                    "id": 2,
                    "title": "The Alchemist",
                    "author": "Paulo Coehlo",
                    "published_date": "1988-07-24",
                    "edition": "25th",
                    "description": "An enchanting story about an Andalusian shepherd boy.",
                    "genre": "Fiction"
                },
                ...  
            ]
        }

        Notes: 
            - The API will return a collection object corresponding to the provided collection_id in JSON format.
            - Each colleciton object contains an `id`, `name`, and list of `books` in that collection.

        Error Handling:
            - 404 Not Found: If this collection is not in the system.

    - Add book to collection

        URL: `/collections/:collection_id`

        Method: POST
        Request Body: 
        {
            "collection_id": 1,
            "book_id": 1
        }
        Response:
        HTTP Status: 201 Created
        {
            "message": "Book successfully added to collection.", 
            "collection_id": 1,
            "book_id": 1
        }

    - Remove book from collection

        URL: `/collections/:collection_id/books/:book_id`

        Method: DELETE
        Response:
        HTTP Status: 204 No Content
        {
            "message": "Book successfully removed from collection."
        }

        Notes:
            - Provide the book_id and collection_id in the URL to specify the book you want to remove from a specific collection.
        
        Error Handling:
            - 404 Not Found: If the specified book is not found in the particular collection.

    - Delete collection

        URL: `/collections/:collection_id`

        Method: DELETE
        Response:
        HTTP Status: 204 No Content
        {
            "message": "Collection successfully deleted."
        }

        Notes:
            - Provide the collection_id in the URL to specify collection you'd like to remove.
        
        Error Handling:
            - 404 Not Found: If the specified collection is not found in the database.











