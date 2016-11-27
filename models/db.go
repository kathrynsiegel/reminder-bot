package models

import "github.com/jinzhu/gorm"

// Database struct contains the gorm struct with a postgres connection.
type Database struct {
	Gorm *gorm.DB
}
