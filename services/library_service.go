package services

import (
	"errors"
	"fmt"
	"sync"

	"example.com/library_management/models"
)

// LibraryManager defines available operations on the library.
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
}

// Library implements LibraryManager.
type Library struct {
	mu      sync.Mutex
	books   map[int]models.Book   // bookID -> Book
	members map[int]models.Member // memberID -> Member
}

// NewLibrary returns an initialized Library.
func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

// AddBook adds a book to the library. If the ID exists it will overwrite.
func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if book.Status == "" {
		book.Status = "Available"
	}
	l.books[book.ID] = book
}

// RemoveBook removes a book by ID if it exists and is not borrowed.
func (l *Library) RemoveBook(bookID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if b.Status == "Borrowed" {
		return errors.New("cannot remove a borrowed book")
	}
	delete(l.books, bookID)
	return nil
}

// BorrowBook allows a member to borrow an available book.
func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, bok := l.books[bookID]
	if !bok {
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}
	member, mok := l.members[memberID]
	if !mok {
		return fmt.Errorf("member with id %d not found", memberID)
	}

	// mark book borrowed and add copy to member.BorrowedBooks
	book.Status = "Borrowed"
	l.books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member
	return nil
}

// ReturnBook allows a member to return a book they borrowed.
func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, bok := l.books[bookID]
	if !bok {
		return errors.New("book not found")
	}
	member, mok := l.members[memberID]
	if !mok {
		return fmt.Errorf("member with id %d not found", memberID)
	}

	// Find the book in the member's borrowed slice
	idx := -1
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("this member did not borrow this book")
	}

	// remove from member borrowed slice
	member.BorrowedBooks = append(member.BorrowedBooks[:idx], member.BorrowedBooks[idx+1:]...)
	l.members[memberID] = member

	// mark book available again
	book.Status = "Available"
	l.books[bookID] = book
	return nil
}

// ListAvailableBooks returns all books with Status == "Available".
func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]models.Book, 0)
	for _, b := range l.books {
		if b.Status == "Available" {
			out = append(out, b)
		}
	}
	return out
}

// ListBorrowedBooks returns all books borrowed by a member.
func (l *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	member, ok := l.members[memberID]
	if !ok {
		return nil, fmt.Errorf("member with id %d not found", memberID)
	}
	return member.BorrowedBooks, nil
}

// Helpers for controller setup:

// AddMember adds or replaces a member in the library.
func (l *Library) AddMember(m models.Member) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.members[m.ID] = m
}

// GetMemberExists checks if member exists.
func (l *Library) GetMemberExists(memberID int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.members[memberID]
	return ok
}
