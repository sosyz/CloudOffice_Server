package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	OpenId  string `gorm:"type:varchar(100);not null;index:idx_openid"`
	Gender  uint   `gorm:"type:short int; default:'0'"`
	Address string `gorm:"type:varchar(100)"`

	// 关联权限组
	AuthId Auth `gorm:"association_foreignkey:ID"`
}

// ChangeStorage 更新用户容量
func (user *User) ChangeStorage(tx *gorm.DB, operator string, size uint64) error {
	return tx.Model(user).Update("storage", gorm.Expr("storage "+operator+" ?", size)).Error
}
