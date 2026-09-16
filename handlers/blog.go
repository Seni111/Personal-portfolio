package handlers

import (
	"log"
	"net/http"

	"portfolio/database"
	"portfolio/models"

	"github.com/gin-gonic/gin"
)

func BlogList(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, title, content, created_at FROM blog_posts ORDER BY created_at DESC")
	if err != nil {
		log.Println("Error fetching blog posts:", err)
		c.HTML(http.StatusInternalServerError, "blog.html", gin.H{
			"title": "Blog",
			"error": "Failed to load posts",
		})
		return
	}
	defer rows.Close()

	var posts []models.BlogPost
	for rows.Next() {
		var p models.BlogPost
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			log.Println("Error scanning post:", err)
			continue
		}
		posts = append(posts, p)
	}

	c.HTML(http.StatusOK, "blog.html", gin.H{
		"title": "Blog",
		"posts": posts,
	})
}

func BlogPost(c *gin.Context) {
	id := c.Param("id")

	var p models.BlogPost
	err := database.DB.QueryRow(
		"SELECT id, title, content, created_at FROM blog_posts WHERE id = $1", id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)

	if err != nil {
		c.HTML(http.StatusNotFound, "post.html", gin.H{
			"title": "Not Found",
			"error": "Post not found",
		})
		return
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"title": p.Title,
		"post":  p,
	})
}
