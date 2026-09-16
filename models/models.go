package models

import "time"

type BlogPost struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}

type Project struct {
	ID          int
	Title       string
	Description string
	URL         string
	ImageURL    string
}

type ContactMessage struct {
	ID        int
	Name      string
	Email     string
	Message   string
	CreatedAt time.Time
}
