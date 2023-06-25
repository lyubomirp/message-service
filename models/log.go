package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model

	Body string
}
