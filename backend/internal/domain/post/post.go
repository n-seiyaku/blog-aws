package post

import "time"

type Post struct {
	ID        string
	Title     string
	Content   string
	ImageKey  string
	AuthorID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
