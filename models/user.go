package models

type User struct {
	UserID    uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Email     string
	Firstname string
	Lastname  string
	Role      string `gorm:"type:enum('Customer','Manager','Sales','Admin');not null"`
}

func (User) TableName() string {
	return "user"
}
