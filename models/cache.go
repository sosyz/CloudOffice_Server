package models

import (
	"time"
)

type Cache struct {
	Key      string    `form:"varchar(128) notnull unique"`
	Value    string    `form:"varchar(1024) notnull"`
	ExpireAt time.Time `form:"datetime notnull"`
}

func (c *Cache) Creat() error {
	tx := DB.Begin()

	if err := tx.Create(c).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (c *Cache) GetValue() (string, error) {
	if err := DB.Where("key = ? AND expireAt > ?", c.Key, time.Now()).First(c).Error; err != nil {
		return "", err
	}
	return c.Value, nil
}

func (c *Cache) Delete() error {
	tx := DB.Begin()

	if err := tx.Where("key = ?", c.Key).Delete(c).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
