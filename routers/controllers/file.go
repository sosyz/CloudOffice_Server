package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/pkg/cos"
	"sonui.cn/cloudprint/pkg/utils"
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
	info.PageNum, err = cos.GetFilePagesNum(record.Path)
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
			"message": "ok",
		})
	}
}

func UploadStart(c *gin.Context) {
	var record models.File
	record.Name = c.PostForm("name")
	record.Path = c.PostForm("path")
	record.Status = models.FileStatusUploadStart

	if record.Path == "" || record.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "option is required",
		})
		return
	}
	record.Fid = utils.SnowFlake.GetId()
	if err := record.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    201,
			"message": "Create file record error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
		})
	}
}
