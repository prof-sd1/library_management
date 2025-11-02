package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"example.com/library_management/models"
	"example.com/library_management/services"
)

// StartConsole runs a simple console interface for the provided library.
// rawLib should be the concrete *services.Library if you want seed helpers to be available.
func StartConsole(lib services.LibraryManager, rawLib interface{}) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nLibrary Menu:")
		fmt.Println("1) Add Book")
		fmt.Println("2) Remove Book")
		fmt.Println("3) Borrow Book")
		fmt.Println("4) Return Book")
		fmt.Println("5) List Available Books")
		fmt.Println("6) List Borrowed Books (by member)")
		fmt.Println("7) Exit")
		fmt.Print("Choose an option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			// Add Book
			fmt.Print("Book ID: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			fmt.Print("Title: ")
			title, _ := reader.ReadString('\n')
			fmt.Print("Author: ")
			author, _ := reader.ReadString('\n')
			book := models.Book{ID: id, Title: strings.TrimSpace(title), Author: strings.TrimSpace(author)}
			lib.AddBook(book)
			fmt.Println("Book added.")

		case "2":
			fmt.Print("Book ID to remove: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			err := lib.RemoveBook(id)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book removed.")
			}

		case "3":
			fmt.Print("Member ID: ")
			midStr, _ := reader.ReadString('\n')
			mid, _ := strconv.Atoi(strings.TrimSpace(midStr))
			fmt.Print("Book ID to borrow: ")
			bidStr, _ := reader.ReadString('\n')
			bid, _ := strconv.Atoi(strings.TrimSpace(bidStr))
			err := lib.BorrowBook(bid, mid)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book borrowed.")
			}

		case "4":
			fmt.Print("Member ID: ")
			midStr, _ := reader.ReadString('\n')
			mid, _ := strconv.Atoi(strings.TrimSpace(midStr))
			fmt.Print("Book ID to return: ")
			bidStr, _ := reader.ReadString('\n')
			bid, _ := strconv.Atoi(strings.TrimSpace(bidStr))
			err := lib.ReturnBook(bid, mid)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book returned.")
			}

		case "5":
			books := lib.ListAvailableBooks()
			if len(books) == 0 {
				fmt.Println("No available books.")
			} else {
				fmt.Println("Available books:")
				for _, b := range books {
					fmt.Printf("ID:%d Title:%s Author:%s\n", b.ID, b.Title, b.Author)
				}
			}

		case "6":
			fmt.Print("Member ID: ")
			midStr, _ := reader.ReadString('\n')
			mid, _ := strconv.Atoi(strings.TrimSpace(midStr))
			books, err := lib.ListBorrowedBooks(mid)
			if err != nil {
				fmt.Println("Error:", err)
			} else if len(books) == 0 {
				fmt.Println("This member has no borrowed books.")
			} else {
				fmt.Printf("Borrowed by member %d:\n", mid)
				for _, b := range books {
					fmt.Printf("ID:%d Title:%s Author:%s\n", b.ID, b.Title, b.Author)
				}
			}

		case "7":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option")
		}
	}
}
