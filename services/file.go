package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/utils"
	"sonui.cn/cloudprint/utils/log"
	"sonui.cn/cloudprint/utils/openxml"
	"time"
)

// SaveFile 保存文件至COS
// name 文件名
// localPath 本地文件路径
// info 文件信息
// size 文件大小
// body 文件内容
func SaveFile(user, name, localPath, info string, size uint64) (*models.File, *log.ErrorInfo) {
	service := s3.New(utils.S3)
	fp, _ := os.Open(localPath)
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
		// 判断是否为docx文件
		if utils.CheckFileTypeBySuffix(name, "docx") && utils.CheckFileTypeByHeader(localPath, []byte("PK")) {
			// 为docx文件，本地解析
			fInfo.PageNum, err = openxml.GetDocxPages(localPath)
		} else {
			fInfo.PageNum, err = utils.GetFilePagesNum(localPath)
		}

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
