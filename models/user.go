package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/utils/log"
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

func GetUserFromRequest(c *gin.Context) (*User, *log.ErrorInfo) {
	//获取cookie
	var openid, session *http.Cookie
	var err error
	openid, err = c.Request.Cookie("openid")
	if err == nil {
		session, err = c.Request.Cookie("session")
	}
	if err != nil || openid.Value == "" || session.Value == "" {
		return nil, log.NewError(301, "no login")
	}

	// 查询登录信息
	var value struct {
		Uid int64 `json:"uid"`
	}

	cache := Cache{
		Key:   session.Value,
		Value: value,
	}

	// 空记录
	if err = cache.Find(); err != nil {
		return nil, log.NewError(301, "no login")
	}

	// 查询用户信息
	user := User{
		Id: value.Uid,
	}
	if err = user.Find(); err != nil {
		return nil, log.NewError(301, "no login")
	}

	return &user, nil
}
