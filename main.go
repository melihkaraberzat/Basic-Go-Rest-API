package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID	   string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Quantity int  `json:"quantity"`
}

var DB *gorm.DB

func connect(){
	db,err := gorm.Open(postgres.Open("postgres://root:root@localhost/postgres?sslmode=disable"),&gorm.Config{})
	if err!=nil{
	log.Fatal("Failed to connect database.")
	}
	DB = db
	db.AutoMigrate(&Book{})
	
}

func addBook(c *gin.Context){
	var book Book
	err := c.BindJSON(&book)
	if err!= nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	}
	result := DB.Create(&book)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func getBooks(c *gin.Context){
	var books []Book
	result := DB.Find(&books)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	c.JSON(http.StatusOK,books)

}

func addBooks(c *gin.Context){
	var payload Book
	c.BindJSON(&payload)
	DB.Create(&payload)
}

func deleteBook(c *gin.Context){
	var book Book
	c.BindJSON(&book)
	DB.Delete(&book)
}

func main() {
	connect()
	router := gin.Default()
	router.POST("/books",addBooks)
	router.GET("/books",getBooks)
	router.DELETE("/books",deleteBook)
	router.Run("localhost:8080")
	
	
	
}