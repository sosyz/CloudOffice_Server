package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	OpenId  string `gorm:"type:varchar(100);not null;index:idx_openid"`
	Gender  uint   `gorm:"type:smallint; default:'0'"`
	Address string `gorm:"type:varchar(100)"`

	// 关联权限组
	// AuthId Auth `gorm:"association_foreignkey:ID"`
}
