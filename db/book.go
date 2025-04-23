package db

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name  string `gorm:"size:255;not null;unique" json:"name"`
	Books []Book `gorm:"foreignKey:AuthorID" json:"books"`
}

type Book struct {
	gorm.Model
	Title    string `gorm:"size:255;not null" json:"title"`
	AuthorID uint   `gorm:"not null" json:"author_id"`
	Author   Author `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"author"`
}

func GetBookByName(title string) (Book, error) {
	var book Book

	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
		return Book{}, fmt.Errorf("failed to connect db: %v", err)
	}

	if err := db.Where("title = ?", title).First(&book).Error; err != nil {
		log.Printf("Failed to get book by name: %v", err)
		return Book{}, fmt.Errorf("failed to get book by name: %v", err)
	}

	return book, nil
}

func GetAllBooks() ([]Book, error) {
	var books []Book

	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
		return []Book{}, fmt.Errorf("failed to connect db: %v", err)
	}

	if err := db.Find(&books).Error; err != nil {
		log.Printf("Failed to get all books: %v", err)
		return []Book{}, fmt.Errorf("failed to get all books: %v", err)
	}

	return books, nil
}

func GetBookByAuthor(author string) ([]Book, error) {
	var books []Book

	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	if err := db.Find(&books, "author = ?", author).Error; err != nil {
		log.Printf("Failed to get book by author: %v", err)
		return []Book{}, fmt.Errorf("failed to get book by author: %v", err)
	}

	return books, nil
}

func DeleteBook(title string) error {
	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	if err := db.Unscoped().Where("title = ?", title).Delete(&Book{}).Error; err != nil {
		log.Printf("Failed to delete book: %v", err)
		return fmt.Errorf("failed to delete book: %v", err)
	}

	return nil
}
