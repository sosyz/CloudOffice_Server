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

	if file, err := services.SaveFile(openid, filename, dst, "", uint64(fileSize)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5,
			"message": "file is required",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":   0,
			"fid":    strconv.FormatInt(file.Fid, 10),
			"status": strconv.Itoa(int(file.Status)),
		})
	}
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
