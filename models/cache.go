package models

import (
	"encoding/json"
	"time"
)

type Cache struct {
	Key      string    `gorm:"type:varchar(128);notnull;unique;"`
	Value    any       `gorm:"-"`
	ExpireAt time.Time `gorm:"type:datetime;notnull"`
	V        string    `gorm:"column:value;type:varchar(10240);notnull;"`
}

// Create 创建缓存
func (c *Cache) Create() error {
	// json编码
	value, err := json.Marshal(c.Value)
	if err != nil {
		return err
	}
	// 写入缓存
	c.V = string(value)
	if err := DB.Create(c).Error; err != nil {
		return err
	}
	return nil
}

// GetValue 查询缓存
func (c *Cache) GetValue() error {
	if err := DB.Where("`key` = ? AND expire_at > ?", c.Key, time.Now()).First(c).Error; err != nil {
		return err
	}
	err := json.Unmarshal([]byte(c.V), c.Value)
	if err != nil {
		return err
	}
	return nil
	//fmt.Println(c.Value)
	//// json解析
	//if err := json.Unmarshal([]byte(c.Value), v); err != nil {
	//	return err
	//}
	//fmt.Printf("%v\n", v)
	//return nil
}

// Delete 删除缓存
func (c *Cache) Delete() error {
	if err := DB.Where("key = ?", c.Key).Delete(Cache{}).Error; err != nil {
		return err
	}
	return nil
}

// ClearOld 清理过期的缓存
func (c *Cache) ClearOld() error {
	if err := DB.Where("expireAt < ?", time.Now()).Delete(Cache{}).Error; err != nil {
		return err
	}
	return nil
}
