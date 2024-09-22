package structures

import "time"

// Status is an enumeration type for the status column
type Status int

const (
	StatusOpen       Status = iota // 0
	StatusInProgress               // 1
	StatusClosed                   // 2
)

type Database struct {
	host     string
	database string
	port     string
	user     string
	password string
	sslmode  string
}

// User model
type Users struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	Username   string    `gorm:"not null" json:"username"`
	Password   string    `gorm:"not null" json:"password"`
	Email      string    `gorm:"not null" json:"email"`
	ProfilePic string    `gorm: "not null"; default:'' json: "profile_pic"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Errors model with added Status and TimeTaken columns
type Errors struct {
	ID          string      `gorm:"primaryKey" json:"id"`
	UserID      int         `gorm:"not null" json:"user_id"`
	Title       string      `gorm:"not null" json:"title"`
	DocFilePath string      `gorm:"not null;default:''" json:"doc_file_path"`
	Image1      string      `gorm:"not null;default:''" json:"image_1"`
	Image2      string      `gorm:"not null;default:''" json:"image_2"`
	Image3      string      `gorm:"not null;default:''" json:"image_3"`
	Image4      string      `gorm:"not null;default:''" json:"image_4"`
	Status      Status      `gorm:"not null;default:0" json:"status"`
	Type        JourneyType `gorm:"not null;default:0" json:"type"`
	TimeTaken   string      `gorm:"not null" json:"time_taken"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Doc model with added Status and TimeTaken columns
// type Docs struct {
// 	ID          string    `gorm:"primaryKey" json:"id"`
// 	UserID      int       `gorm:"not null" json:"user_id"`
// 	Title       string    `gorm:"not null" json:"title"`
// 	DocFilePath string    `gorm:"not null;default:''" json:"doc_file_path"`
// 	Image1      string    `gorm:"not null;default:''" json:"image_1"`
// 	Image2      string    `gorm:"not null;default:''" json:"image_2"`
// 	Image3      string    `gorm:"not null;default:''" json:"image_3"`
// 	Image4      string    `gorm:"not null;default:''" json:"image_4"`
// 	Status      Status    `gorm:"not null;default:0" json:"status"`
// 	TimeTaken   string    `gorm:"not null" json:"time_taken"`
// 	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
// 	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
// }

// UsefulType is an enumeration type for the Useful column
type UsefulType int

const (
	NotUseful UsefulType = iota // 0
	Useful                      // 1
)

// UsefulType is an enumeration type for the Useful column
type JourneyType int

const (
	ErrorJourney JourneyType = iota // 0
	DocJourney                      // 1
)

// UserActions model with updated Useful column as an enum
type UserActions struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	UserID      int        `gorm:"not null" json:"user_id"`
	ErrorID     string     `gorm:"not null" json:"error_id"`
	TextContent string     `gorm:"not null" json:"text_content"`
	Type        string     `gorm:"not null" json:"type"`
	Useful      UsefulType `gorm:"not null" json:"useful"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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
