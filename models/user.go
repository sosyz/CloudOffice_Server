package models

import (
	"time"
)

type User struct {
	Id         int64  `json:"id" gorm:"primary_key"`
	Name       string `gorm:"type:varchar(100);not null"`
	OpenId     string `gorm:"type:varchar(100);not null;index:idx_openid"`
	Gender     uint   `gorm:"type:smallint; default:'0'"`
	Address    string `gorm:"type:varchar(100)"`
	Phone      string `gorm:"type:char(24)"`
	CreateTime time.Time
	UpdateTime time.Time
	Group      uint
	// 关联权限组
	// AuthId Auth `gorm:"association_foreignkey:ID"`
}

// Find 查找用户
func (u *User) Find() error {
	return DB.Where("open_id = ?", u.OpenId).First(u).Error
}

// Create 创建用户
func (u *User) Create() error {
	u.CreateTime = time.Now()
	return DB.Create(u).Error
}

func (u *User) Save() error {
	u.UpdateTime = time.Now()
	return DB.Model(u).Updates(u).Error
}
