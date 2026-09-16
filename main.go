package main

import (
	"log"

	"portfolio/database"
	"portfolio/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", handlers.Home)
	r.GET("/projects", handlers.Projects)
	r.GET("/blog", handlers.BlogList)
	r.GET("/blog/:id", handlers.BlogPost)
	r.GET("/contact", handlers.ContactPage)
	r.POST("/contact", handlers.ContactSubmit)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
