package main

import (
	"fmt"

	"example.com/library_management/controllers"
	"example.com/library_management/models"
	"example.com/library_management/services"
)

func main() {
	lib := services.NewLibrary()

	// seed some members
	lib.AddMember(models.Member{ID: 1, Name: "Alice"})
	lib.AddMember(models.Member{ID: 2, Name: "Bob"})

	// seed some books
	lib.AddBook(models.Book{ID: 101, Title: "The Go Programming Language", Author: "Alan A. A. Donovan"})
	lib.AddBook(models.Book{ID: 102, Title: "Clean Code", Author: "Robert C. Martin"})
	lib.AddBook(models.Book{ID: 103, Title: "Introduction to Algorithms", Author: "Cormen et al."})

	fmt.Println("Welcome to the Console Library Management System")
	controllers.StartConsole(lib, lib)
}
