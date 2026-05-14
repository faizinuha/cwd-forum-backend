package model

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Notification struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ThreadId  *uint      `json:"thread_id,omitempty"`
	PostId    *uint      `json:"post_id,omitempty"`
	UserId    uint       `json:"user_id"`
	Type      string     `json:"type"`    // e.g., "mention", "reply", "upvote"
	Payload   string     `json:"payload"` // JSON string with additional data (e.g., thread ID, post ID)
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`

	// Associations
	User   User    `gorm:"foreignKey:UserId" json:"user"`
	Thread *Thread `gorm:"foreignKey:ThreadId" json:"thread,omitempty"`
	Post   *Post   `gorm:"foreignKey:PostId" json:"post,omitempty"`
}

type User struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type Thread struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Content  string    `json:"content"`
}

type Post struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Content  string    `json:"content"`
	AuthorId uint      `json:"author_id"`
}

var DB_DRIVER = "sqlite3"
var DB_SOURCE = "database/forum.db"
