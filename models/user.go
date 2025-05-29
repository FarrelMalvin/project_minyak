package models

type User struct {
	UserID    uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueindex;not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"uniqueindex;not null"`
	Firstname string
	Lastname  string
	Role      string `gorm:"type:enum('Customer','Manager','Sales','Admin');not null;index"`
}

func (User) TableName() string {
	return "user"
}
