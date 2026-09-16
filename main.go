package main

import (
	"html/template"
	"log"

	"portfolio/database"
	"portfolio/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type templateRenderer struct {
	templates map[string]*template.Template
}

func (t *templateRenderer) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: t.templates[name],
		Name:     "layout",
		Data:     data,
	}
}

func loadTemplates() *templateRenderer {
	templates := make(map[string]*template.Template)
	pages := []string{"home.html", "blog.html", "post.html", "projects.html", "contact.html"}
	for _, page := range pages {
		templates[page] = template.Must(
			template.ParseFiles("templates/layout.html", "templates/"+page),
		)
	}
	return &templateRenderer{templates: templates}
}

func main() {
	database.Connect()

	r := gin.Default()

	r.HTMLRender = loadTemplates()
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
