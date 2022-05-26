package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	OpenId  string `gorm:"type:varchar(100);not null;index:idx_openid"`
	Gender  uint   `gorm:"type:smallint; default:'0'"`
	Address string `gorm:"type:varchar(100)"`
	Phone   string `gorm:"type:char(24)"`
	Group   uint
	// 关联权限组
	// AuthId Auth `gorm:"association_foreignkey:ID"`
}

// Get 查找用户
func (u *User) Find() error {
	return DB.Where("open_id = ?", u.OpenId).First(u).Error
}

// Create 创建用户
func (u *User) Create() error {
	return DB.Create(u).Error
}

func (u *User) Save() error {
	return DB.Model(u).Updates(u).Error
}
