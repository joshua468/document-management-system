package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Email    string
}

func (u *User) CheckPassword(password string) bool {
	// Implement bcrypt comparison in production
	return u.Password == password
}
