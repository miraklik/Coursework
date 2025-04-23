package main

import (
	"bookstore/db"
	"bookstore/handlers"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	return db.ConnectDB()
}

func main() {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Bookstore")

	titleEntry := widget.NewEntry()
	authorEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	titleBookEntry := widget.NewEntry()
	authorNameEntry := widget.NewEntry()
	titleForDeleteBookEntry := widget.NewEntry()

	createBookButton := widget.NewButton("Create Book", func() {
		title := titleEntry.Text
		author := authorEntry.Text

		book, err := handlers.CreateBook(db, title, author)
		if err != nil {
			log.Printf("Failed to create book: %v", err)
			return
		} else {
			log.Printf("Book created: %v", book)
		}
	})

	createGetAllBookButton := widget.NewButton("Get All Books", func() {
		books, err := handlers.GetAllBooks(db)
		if err != nil {
			log.Printf("Failed to get all books: %v", err)
			return
		} else {
			log.Printf("All books: %v", books)
		}
	})

	createGetBookByNameButton := widget.NewButton("Get Book By Name", func() {
		title := titleBookEntry.Text

		books, err := handlers.GetBookByName(db, title)
		if err != nil {
			log.Printf("Failed to get book by name: %v", err)
			return
		} else {
			log.Printf("Books by name: %v", books)
		}
	})

	createGetBookByAuthorButton := widget.NewButton("Get Book By Author", func() {
		authorName := authorNameEntry.Text

		books, err := handlers.GetBookByAuthor(db, authorName)
		if err != nil {
			log.Printf("Failed to get book by author: %v", err)
			return
		} else {
			log.Printf("Books by author: %v", books)
		}
	})

	createDeleteBookButton := widget.NewButton("Delete Book", func() {
		title := titleForDeleteBookEntry.Text

		err := handlers.DeleteBook(db, title)
		if err != nil {
			log.Printf("Failed to delete book: %v", err)
			return
		} else {
			log.Println("Book deleted")
		}
	})

	createAuthorButton := widget.NewButton("Create Author", func() {
		name := nameEntry.Text

		author, err := handlers.CreateAuthor(db, name)
		if err != nil {
			log.Printf("Failed to create author: %v", err)
			return
		} else {
			log.Printf("Author created: %v", author)
		}
	})

	form := container.NewVBox(
		widget.NewLabel("Create Book:"),
		titleEntry,
		authorEntry,
		createBookButton,
		widget.NewLabel("Create Author:"),
		nameEntry,
		createAuthorButton,
		widget.NewLabel("Get All Books:"),
		createGetAllBookButton,
		widget.NewLabel("Get Book By Name:"),
		titleBookEntry,
		createGetBookByNameButton,
		widget.NewLabel("Get Book By Author:"),
		authorNameEntry,
		createGetBookByAuthorButton,
		widget.NewLabel("Delete Book:"),
		titleForDeleteBookEntry,
		createDeleteBookButton,
	)

	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}
