package models

import "time"

type Log struct {
	Id int64 `json:"id" gorm:"primary_key"`
	// 操作者
	User int64 `json:"user" gorm:"type:bigint;"`
	// 操作类型
	Type int `json:"type" gorm:"type:int;"`
	// 操作时间
	Time time.Time `json:"time" gorm:"type:datetime;"`
	// 操作内容
	Content string `json:"content" gorm:"type:varchar(100);"`
}

func (l *Log) Create() error {
	return DB.Create(l).Error
}
