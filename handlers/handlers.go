package handlers

import (
	"bookstore/db"
	"fmt"

	"gorm.io/gorm"
)

func CreateBook(database *gorm.DB, title, author string) (*db.Book, error) {
	var authorRecord db.Author
	// Ищем автора по имени
	if err := database.Where("name = ?", author).First(&authorRecord).Error; err != nil {
		// Ошибка при поиске автора
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("author not found")
		}
		return nil, fmt.Errorf("failed to get author: %v", err)
	}

	// Создаем книгу с найденным author_id
	book := db.Book{Title: title, AuthorID: authorRecord.ID}
	if err := database.Create(&book).Error; err != nil {
		return nil, fmt.Errorf("failed to create book: %v", err)
	}

	return &book, nil
}

func GetAllBooks(database *gorm.DB) (*[]db.Book, error) {
	var books []db.Book
	if err := database.Find(&books).Error; err != nil {
		return nil, fmt.Errorf("failed to get all books: %v", err)
	}

	return &books, nil
}

func GetBookByName(database *gorm.DB, title string) (*db.Book, error) {
	var book db.Book
	if err := database.Where("title = ?", title).First(&book).Error; err != nil {
		return nil, fmt.Errorf("failed to get book by name: %v", err)
	}

	return &book, nil
}

func GetBookByAuthor(database *gorm.DB, author string) (*[]db.Book, error) {
	var books []db.Book
	if err := database.Where("author_id IN (SELECT id FROM authors WHERE name = ?)", author).Find(&books).Error; err != nil {
		return nil, fmt.Errorf("failed to get book by author: %v", err)
	}

	return &books, nil
}

func DeleteBook(database *gorm.DB, title string) error {
	if err := database.Where("title = ?", title).Delete(&db.Book{}).Error; err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}

	return nil
}

func CreateAuthor(database *gorm.DB, name string) (*db.Author, error) {
	author := db.Author{Name: name}
	if err := database.Create(&author).Error; err != nil {
		return nil, fmt.Errorf("failed to create author: %v", err)
	}

	return &author, nil
}
