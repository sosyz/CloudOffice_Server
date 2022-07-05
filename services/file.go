package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/utils"
	"sonui.cn/cloudprint/utils/log"
	"time"
)

// SaveFile 保存文件至COS
// name 文件名
// path 本地文件路径
// info 文件信息
// size 文件大小
// body 文件内容
func SaveFile(user, name, path, info string, size uint64) (*models.File, *log.ErrorInfo) {
	service := s3.New(utils.S3)
	fp, _ := os.Open(path)
	defer utils.FileClose(fp)

	withContext, err := service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(utils.Config.QCloud.Bucket),
		Key:    aws.String(user + "/" + name),
		Body:   fp,
	})

	if err != nil {
		return nil, log.NewError(103, err.Error())
	} else {
		file := models.File{
			Fid:        utils.FileSF.GetId(),
			Name:       name,
			Path:       user + "/" + name,
			Size:       size,
			Status:     0,
			Info:       info,
			ETage:      *withContext.ETag,
			CreateTime: time.Time{},
			UpdateTime: time.Time{},
		}
		var fInfo = models.FileInfo{}
		fInfo.PageNum, err = utils.GetFilePagesNum(path)

		if err != nil {
			return nil, log.NewError(101, err.Error())
		}

		file.Info = fInfo.ToJson()
		file.Status = models.FileStatusUploadCompacter

		if err := file.Create(); err != nil {
			return nil, log.NewError(201, err.Error())
		}
		return &file, nil
	}
}

// GetFile 获取文件
// fid 文件id
func GetFile(fid int64) (*s3.GetObjectOutput, *log.ErrorInfo) {
	file := models.File{
		Fid: fid,
	}
	err := file.Find()
	if err != nil {
		return nil, log.NewError(203, "not found this record")
	}

	service := s3.New(utils.S3)

	object, err := service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(utils.Config.QCloud.Bucket),
		Key:    aws.String(file.Path),
	})
	if err != nil {
		return nil, log.NewError(103, err.Error())
	}

	return object, nil
}
