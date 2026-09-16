package handlers

import (
	"log"
	"net/http"

	"portfolio/database"
	"portfolio/models"

	"github.com/gin-gonic/gin"
)

func Projects(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, title, description, COALESCE(url, ''), COALESCE(image_url, '') FROM projects ORDER BY id DESC")
	if err != nil {
		log.Println("Error fetching projects:", err)
		c.HTML(http.StatusInternalServerError, "projects.html", gin.H{
			"title": "Projects",
			"error": "Failed to load projects",
		})
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.URL, &p.ImageURL); err != nil {
			log.Println("Error scanning project:", err)
			continue
		}
		projects = append(projects, p)
	}

	c.HTML(http.StatusOK, "projects.html", gin.H{
		"title":    "Projects",
		"projects": projects,
	})
}
