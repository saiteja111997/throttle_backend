package structures

import "time"

type Database struct {
	host     string
	database string
	port     string
	user     string
	password string
	sslmode  string
}

// User model
type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Error model
type Errors struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	Title     string    `gorm:"not null" json:"title"`
	FilePath  string    `gorm:"not null default:''" json:"file_path"`
	Image1    string    `gorm:"not null default:''" json:"image_1"`
	Image2    string    `gorm:"not null default:''" json:"image_2"`
	Image3    string    `gorm:"not null default:''" json:"image_3"`
	Image4    string    `gorm:"not null default:''" json:"image_4"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// UserActions model
type UserActions struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	UserID      int       `gorm:"not null" json:"user_id"`
	ErrorID     string    `gorm:"not null" json:"error_id"`
	TextContent string    `gorm:"not null" json:"text_content"`
	Type        string    `gorm:"not null" json:"type"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// CREATE TABLE users (
//     id SERIAL PRIMARY KEY,
//     username VARCHAR(255) NOT NULL,
//     password VARCHAR(255) NOT NULL,
//     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
// );

// CREATE TABLE errors (
//     id VARCHAR(36) PRIMARY KEY,
//     user_id INTEGER REFERENCES users(id),
//     title VARCHAR(255) NOT NULL,
//     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
// );

// CREATE TABLE user_actions (
//     id SERIAL PRIMARY KEY,
//     user_id INTEGER REFERENCES users(id),
//     error_id VARCHAR(36) REFERENCES errors(id),
//     text_content TEXT NOT NULL,
//     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
// );

// ALTER TABLE errors
// ADD COLUMN FilePath VARCHAR NOT NULL DEFAULT '',
// ADD COLUMN Image1 VARCHAR NOT NULL DEFAULT '',
// ADD COLUMN Image2 VARCHAR NOT NULL DEFAULT '',
// ADD COLUMN Image3 VARCHAR NOT NULL DEFAULT '';
// ADD COLUMN Image4 VARCHAR NOT NULL DEFAULT '';
