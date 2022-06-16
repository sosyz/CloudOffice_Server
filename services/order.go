package services

import (
	"errors"
	"fmt"
	"sonui.cn/cloudprint/models"
)

type OrderOverview struct {
	OrderId      int64    `json:"order_id"`
	FileNameList []string `json:"file_name_list"`
	Status       int      `json:"status"`
}

func GetOrderOverviewList(userId string) ([]OrderOverview, error) {
	list, err := models.OrderList(userId)
	if err != nil {
		return nil, errors.New(err.Error())
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

func OrderInfo(orderId int64) (models.Order, error) {
	order := models.Order{
		ID: orderId,
	}

	if err := order.Find(); err != nil {
		return models.Order{}, fmt.Errorf("get order info error: %s", err.Error())
	}
	return order, nil
}
