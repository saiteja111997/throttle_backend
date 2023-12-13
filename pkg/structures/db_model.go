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
	ID        int       `gorm:"primaryKey"`
	Username  string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Error model
type Error struct {
	ID        string    `gorm:"primaryKey"`
	UserID    int       `gorm:"not null"`
	Title     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
