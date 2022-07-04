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
	UpdateTime time.Time `gorm:"not null"`
	ETage      string    `gorm:"type:char(36)"`
}

// FileInfo 文件信息
type FileInfo struct {
	PageNum   int    `json:"pageNum"`
	Type      int    `json:"type"`      // 纸张类型
	Copies    int    `json:"copies"`    // 份数
	Colour    bool   `json:"colour"`    // 是否彩色
	PrintPage string `json:"printPage"` // 打印范围
}

// ToJson 序列化
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
	file.CreateTime = time.Now()
	file.UpdateTime = time.Now()

	if err := DB.Create(file).Error; err != nil {
		return err
	}
	return nil
}

// Save 保存变动
func (file *File) Save() error {
	file.UpdateTime = time.Now()
	if err := DB.Model(&file).Update(file).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除记录
func (file *File) Delete() error {
	if err := DB.Delete(file).Error; err != nil {
		return err
	}
	return nil
}

// Find 查找记录
func (file *File) Find() error {
	return DB.Model(file).Where(file).First(file).Error
}

// Exist 记录是否存在
func (file *File) Exist() (bool, error) {
	var tmp File
	if err := DB.Model(file).Where(file).First(&tmp).Error; err != nil {
		return false, err
	}
	return true, nil
}

// FileList 获取用户上传的文件列表
func FileList(userID string) ([]File, error) {
	var files []File
	if err := DB.Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
