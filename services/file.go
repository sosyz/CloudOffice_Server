package services

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/utils"
	"time"
)

// SaveFile 保存文件至COS
// name 文件名
// path 本地文件路径
// info 文件信息
// size 文件大小
// body 文件内容
func SaveFile(user, name, path, info string, size uint64) error {
	file := models.File{
		Fid:        utils.FileSF.GetId(),
		Name:       name,
		Path:       user + "/" + name,
		Size:       size,
		Status:     0,
		Info:       info,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}
	service := s3.New(utils.S3)
	fp, _ := os.Open(path)
	defer utils.FileClose(fp)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	withContext, err := service.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(utils.Config.QCloud.Bucket),
		Key:    aws.String(user + "/" + name),
		Body:   fp,
	})
	if err != nil {
		return err
	} else {
		file.ETage = *withContext.ETag
		if err := file.Create(); err != nil {
			return err
		}
	}
	return nil
}

// GetFile 获取文件
// fid 文件id
func GetFile(fid int64) (s3.GetObjectOutput, error) {
	file := models.File{
		Fid: fid,
	}
	err := file.Find()
	if err != nil {
		return s3.GetObjectOutput{}, nil
	}

	service := s3.New(utils.S3)

	object, err := service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(utils.Config.QCloud.Bucket),
		Key:    aws.String(file.Path),
	})
	if err != nil {
		return s3.GetObjectOutput{}, err
	}

	return *object, nil
}
