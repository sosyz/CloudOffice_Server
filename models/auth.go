package models

import "github.com/jinzhu/gorm"

type Auth struct {
	gorm.Model
	Name   string `gorm:"type:varchar(100);"`
	Values uint32 `gorm:"type:int(11);"`
}
