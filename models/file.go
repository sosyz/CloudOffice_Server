package models

import (
	"encoding/json"
	"time"
)

// File 文件
type File struct {
	// 表字段
	Fid        int64     `gorm:"primary_key;"`
	Name       string    `gorm:"type:varchar(100);not null"`
	Path       string    `gorm:"type:varchar(255);not null"`
	Size       uint64    `gorm:"type:bigint;not null"`
	Status     uint8     `gorm:"type:tinyint;default:0"`
	Info       string    `gorm:"type:varchar(255);default:'{}'"`
	CreateTime time.Time `gorm:"not null"`
}

// FileInfo 文件信息
type FileInfo struct {
	PageNum   int    `json:"pageNum"`
	Type      int    `json:"type"`      // 纸张类型
	Copies    int    `json:"copies"`    // 份数
	Colour    bool   `json:"colour"`    // 是否彩色
	PrintPage string `json:"printPage"` // 打印范围
}

func (f *FileInfo) ToJson() string {
	// json编码
	value, _ := json.Marshal(f)
	return string(value)
}

const (
	// FileStatusUploadStart 上传中
	FileStatusUploadStart = iota
	// FileStatusUploadCompacter 上传完成
	FileStatusUploadCompacter
	// FileStatusInfoOk 信息获取完成
	FileStatusInfoOk
	// FileStatusPrinting 打印中
	FileStatusPrinting
	// FileStatusPrinted 打印完成
	FileStatusPrinted
	// FileStatusError 打印错误
	FileStatusError
)

// Create 创建文件记录
func (file *File) Create() error {
	if err := DB.Create(file).Error; err != nil {
		return err
	}
	return nil
}

// Save 保存信息
func (file *File) Save() error {
	if err := DB.Model(&file).Update(file).Error; err != nil {
		return err
	}
	return nil
}

func (file *File) Delete() error {
	if err := DB.Delete(file).Error; err != nil {
		return err
	}
	return nil
}

func (file *File) Find() error {
	return DB.Model(file).Where(file).First(file).Error
}

func (file *File) List() ([]File, error) {
	var files []File
	err := DB.Model(file).Where(file).Find(&files).Error
	return files, err
}
