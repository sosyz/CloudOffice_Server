package services

import (
	"encoding/json"
	"fmt"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/utils"
	"sonui.cn/cloudprint/utils/log"
	"time"
)

type OrderOverview struct {
	OrderId      int64    `json:"order_id"`
	FileNameList []string `json:"file_name_list"`
	Status       int      `json:"status"`
}

func GetOrderOverviewList(userId string) ([]OrderOverview, *log.ErrorInfo) {
	list, err := models.OrderList(userId)
	if err != nil {
		return nil, log.NewError(5, err.Error())
	}
	// 对订单信息进行包装脱敏
	var ret []OrderOverview
	for _, v := range list {
		var fileNameList []string
		for _, v2 := range v.FileList {
			file := models.File{
				Fid: v2,
			}
			err := file.Find()
			if err == nil {
				fileNameList = append(fileNameList, file.Name)
			}
		}

		orderOverview := OrderOverview{
			OrderId:      v.ID,
			FileNameList: fileNameList,
			Status:       v.Status,
		}
		ret = append(ret, orderOverview)
	}
	return ret, nil
}

func OrderInfo(orderId int64) (*models.Order, *log.ErrorInfo) {
	order := models.Order{
		ID: orderId,
	}
	err := order.Find()
	if err != nil {
		return nil, log.NewError(5, err.Error())
	}
	return &order, nil
}

func OrderCreate(userId int64, fileIdList []int64) (*models.Order, *log.ErrorInfo) {
	ans := 0

	// 检查文件是否存在 存在计算页数和
	for _, fid := range fileIdList {
		tFile := models.File{Fid: fid}
		if ok, err := tFile.Exist(); !ok || err != nil {
			return nil, log.NewError(203, fmt.Sprintf("fid[%v] is not exist", fid)) // 只要有一个文件不存在，就返回
		} else {
			var info models.FileInfo
			err = json.Unmarshal([]byte(tFile.Info), &info)
			if err != nil {
				return nil, log.NewError(203, fmt.Sprintf("fid[%v] is not exist", fid)) // 只要有一个文件信息错误
			}
			log.Debug("info", fmt.Sprintf("fid: %v, PageNum: %+v, info:%v\n", tFile.Fid, info.PageNum, tFile.Info))
			ans += info.PageNum
		}
	}

	// 生成订单
	order := models.Order{
		ID:        utils.OrderSF.GetId(),
		FileList:  fileIdList,
		Status:    models.OrderStatusWaitPay,
		UserID:    userId,
		TotalFee:  ans * 30,
		CreatedAt: time.Now(),
	}

	// 写入数据库
	if err := order.Create(); err != nil {
		return nil, log.NewError(201, err.Error())
	}

	return &order, nil
}
