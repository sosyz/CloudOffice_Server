package services

import (
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/utils/log"
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
