package handlers

import (
	"bookstore/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) CreateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Title  string `json:"title"`
			Author string `json:"author"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var author db.Author
		if err := s.db.Where("name=?", req.Author).First(&author).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if req.Author == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Author is required"})
			return
		}

		book := db.Book{Title: req.Title, AuthorID: author.ID}
		if err := s.db.Create(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"book": book})
	}
}

func (s *Server) GetAllBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := db.GetAllBooks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"books": books})
	}
}

func (s *Server) GetBookByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Title string `json:"title"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		book, err := db.GetBookByName(req.Title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"book": book})
	}
}

func (s *Server) GetBookByAuthor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Author string `json:"author"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		books, err := db.GetBookByAuthor(req.Author)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"books": books})
	}
}

func (s *Server) DeleteBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Title string `json:"title"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		if err := db.DeleteBook(req.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
	}
}

func (s *Server) CreateAuthor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		author := db.Author{Name: req.Name}
		if err := s.db.Create(&author).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"author": author})
	}
}
