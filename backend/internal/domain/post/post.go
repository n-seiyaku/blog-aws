package post

import "time"

type Post struct {
	ID        string
	Title     string
	AuthorID  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
