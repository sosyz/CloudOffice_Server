package models

import "sonui.cn/cloudprint/pkg/utils"

// File 文件
type File struct {
	// 表字段
	Fid    int64  `gorm:"primary_key;"`
	Name   string `gorm:"type:varchar(100);not null"`
	Path   string `gorm:"type:varchar(255);not null"`
	Size   uint64 `gorm:"type:bigint;not null"`
	Status uint8  `gorm:"type:tinyint;default:0"`
	// 外键连接到User.Id
	UserId  string `gorm:"column:user_id;size:32;not null;index:idx_user_id" json:"user_id"`
	ForUser User   `gorm:"association_foreignkey:ID"`
}

func (file *File) Create() error {
	tx := DB.Begin()

	if err := tx.Create(file).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (file *File) CreatFid() {
	// 生成节点实例
	node, err := utils.NewWorker(1)
	if err != nil {
		panic(err)
	}
	file.Fid = node.GetId()
}

func (file *File) Update() error {
	tx := DB.Begin()

	if err := tx.Model(&file).Update(file).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (file *File) Delete() error {
	tx := DB.Begin()

	if err := tx.Model(&file).Delete(file).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (file *File) Find() error {
	return DB.Model(file).Where(file).First(file).Error
}

func (file *File) List() ([]File, error) {
	var files []File
	err := DB.Model(file).Where(file).Find(&files).Error
	return files, err
}
