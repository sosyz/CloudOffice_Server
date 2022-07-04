package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/services"
	"sonui.cn/cloudprint/utils"
	"strconv"
)

func UploadComplete(c *gin.Context) {
	var record models.File
	var err error
	record.Fid, err = strconv.ParseInt(c.PostForm("fid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "fid is error",
		})
		return
	}
	if err = record.Find(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    202,
			"message": "fid is error",
		})
		return
	}

	var info models.FileInfo
	info.PageNum, err = utils.GetFilePagesNum(record.Path + record.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    107,
			"message": "get file pageNum error",
		})
		return
	}
	record.Info = info.ToJson()
	record.Status = models.FileStatusUploadCompacter
	err = record.Save()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    201,
			"message": "Save file record failed",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"pageNum": strconv.Itoa(info.PageNum),
		})
	}
}

func UploadStart(c *gin.Context) {
	var record models.File
	record.Name = c.PostForm("name")
	record.Path = c.PostForm("openid")
	record.Status = models.FileStatusUploadStart

	if record.Path == "" || record.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "option is required",
		})
		return
	}
	// 判断path是否以/结尾
	if record.Path[len(record.Path)-1:] != "/" {
		record.Path = record.Path + "/"
	}

	record.Fid = utils.FileSF.GetId()
	if err := record.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    201,
			"message": "Create file record error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"fid":  strconv.FormatInt(record.Fid, 10),
		})
	}
}

func FileUpload(c *gin.Context) {
	openid, _ := c.GetPostForm("openid")

	// form-data
	// 获取文件参数
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "file is required",
		})
		return
	}

	// 获取文件名
	filename := formFile.Filename
	// 获取文件大小
	fileSize := formFile.Size
	// 写到本地
	fName, _ := uuid.NewUUID()
	dst := utils.Config.Run.Temp + fName.String()
	fmt.Printf("%s\n", dst)
	if err := c.SaveUploadedFile(formFile, dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4,
			"message": "file is required",
		})
		return
	}

	if err := services.SaveFile(openid, filename, dst, "", uint64(fileSize)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5,
			"message": "file is required",
		})
		return
	}
	fmt.Printf("name: %v, size: %v", filename, fileSize)
}

func FileDownload(c *gin.Context) {
	fid := c.Query("fid")
	fidInt, _ := strconv.ParseInt(fid, 10, 64)
	var record models.File
	record.Fid = fidInt
	if err := record.Find(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    202,
			"message": "fid is error",
		})
		return
	}
	file, err := services.GetFile(fidInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    202,
			"message": "fid is error",
		})
	}
	c.DataFromReader(http.StatusOK, *file.ContentLength, *file.ContentType, file.Body, nil)
}
