package main

import (
	"net/http"
	"errors"
	"strings"
	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

// Get all books
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

// Get book by ID
func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// Checkout a book
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

// Return a book
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

// Get book by ID helper function
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

// Create a new book
func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// Delete a book
func deleteBook(c *gin.Context) {
	id := c.Param("id")

	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Book deleted successfully."})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
}

// Update book details
func updateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook book

	if err := c.BindJSON(&updatedBook); err != nil {
		return
	}

	for i, b := range books {
		if b.ID == id {
			books[i] = updatedBook
			c.IndentedJSON(http.StatusOK, updatedBook)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
}

// Get books with available stock
func getAvailableBooks(c *gin.Context) {
	var availableBooks []book

	for _, b := range books {
		if b.Quantity > 0 {
			availableBooks = append(availableBooks, b)
		}
	}

	c.IndentedJSON(http.StatusOK, availableBooks)
}

// Search books by title or author
func searchBooks(c *gin.Context) {
	query, exists := c.GetQuery("query")

	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing search query."})
		return
	}

	var results []book
	query = strings.ToLower(query)

	for _, b := range books {
		if strings.Contains(strings.ToLower(b.Title), query) || strings.Contains(strings.ToLower(b.Author), query) {
			results = append(results, b)
		}
	}

	if len(results) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No books found."})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.DELETE("/books/:id", deleteBook)
	router.PUT("/books/:id", updateBook)
	router.GET("/books/available", getAvailableBooks)
	router.GET("/books/search", searchBooks)

	router.Run("localhost:8080")
}
