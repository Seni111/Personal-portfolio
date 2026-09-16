package handlers

import (
	"log"
	"net/http"

	"portfolio/database"

	"github.com/gin-gonic/gin"
)

func ContactPage(c *gin.Context) {
	c.HTML(http.StatusOK, "contact.html", gin.H{
		"title": "Contact",
	})
}

func ContactSubmit(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	message := c.PostForm("message")

	if name == "" || email == "" || message == "" {
		c.HTML(http.StatusBadRequest, "contact.html", gin.H{
			"title": "Contact",
			"error": "All fields are required",
		})
		return
	}

	_, err := database.DB.Exec(
		"INSERT INTO contact_messages (name, email, message) VALUES ($1, $2, $3)",
		name, email, message,
	)
	if err != nil {
		log.Println("Error saving message:", err)
		c.HTML(http.StatusInternalServerError, "contact.html", gin.H{
			"title": "Contact",
			"error": "Failed to send message",
		})
		return
	}

	c.HTML(http.StatusOK, "contact.html", gin.H{
		"title":   "Contact",
		"success": "Message sent successfully!",
	})
}
