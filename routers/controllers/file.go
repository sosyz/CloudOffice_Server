package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
)

const (
	// Init 上传初始化
	Init = iota

	// Complete 上传完成
	Complete

	// FileInfo 文件信息获取完成
	FileInfo

	// OrderMember 订单成员
	OrderMember

	// WaitPrint 打印等待
	WaitPrint

	// Printing 打印中
	Printing

	// PrintComplete 打印完成
	PrintComplete

	// PrintError 打印错误
	PrintError
)

func UploadComplete(c *gin.Context) {
	var record models.File
	if err := c.Bind(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if record.Path == "" || record.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path or name are required"})
		return
	}

	record.Status = Complete
	err := record.Update()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    100,
			"message": "Create file record failed",
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
	if err := c.Bind(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if record.Path == "" || record.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path or name are required"})
		return
	}
	record.CreatFid()
	record.Status = Init
	err := record.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    201,
			"message": "Create file record failed",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
		})
	}
}
