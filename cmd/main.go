package main

import (
	"bookstore/db"
	"bookstore/handlers"
	"log"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	server := handlers.NewServer(db)

	r.POST("/book", server.CreateBook())
	r.POST("/author", server.CreateAuthor())
	r.GET("/book", server.GetAllBooks())
	r.DELETE("/book", server.DeleteBook())
	r.GET("/:name", server.GetBookByName())
	r.GET("/author/:author", server.GetBookByAuthor())

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
