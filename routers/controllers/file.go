package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
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
