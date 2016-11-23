package models

import "github.com/jinzhu/gorm"

// User represents a row in public.users
type User struct {
	gorm.Model
	FacebookID string
}

// FindUserOrCreate finds a user by FB ID. If user isn't present, then
// create a new user in the DB.
func (db *Database) FindUserOrCreate(facebookID string) *User {
	var user User
	db.Gorm.FirstOrCreate(&user, User{FacebookID: facebookID})
	return &user
}
