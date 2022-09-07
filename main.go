package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "booktitle1", Author: "zakky1", Quantity: 1},
	{ID: "2", Title: "booktitle2", Author: "zakky2", Quantity: 1},
	{ID: "3", Title: "booktitle3", Author: "zakky3", Quantity: 1},
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

// ROUTE HANDLERs

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id") // /books/1
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	// get query parameter
	id, ok := c.GetQuery("id")

	// check if id query parameter is exist
	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}

	book, err := getBookById(id)

	// check if book found
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	// check whether book is not out of stock
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	// decrement quantity when checkout is success
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	// get query parameter
	id, ok := c.GetQuery("id")

	// check if id query parameter is exist
	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}

	book, err := getBookById(id)

	// check if book found
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var new_book book

	// get json and create new book
	err := c.BindJSON(&new_book)
	fmt.Println(err)
	if err != nil {
		return
	}

	fmt.Println(new_book.Title)

	books = append(books, new_book)
	c.IndentedJSON(http.StatusCreated, new_book)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
