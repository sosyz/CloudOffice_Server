package models

import (
	"fmt"
	"testing"
	"time"
)

func TestOrder_Create(t *testing.T) {
	order := Order{
		ID:       1,
		UserID:   "test",
		FileList: []int64{1, 2, 3},
		Status:   OrderStatusWaitPay,
		TotalFee: 30,
	}
	err := order.Create()
	if err != nil {
		t.Error(err)
	}
}

func TestOrder_Find(t *testing.T) {
	order := Order{
		ID: 1,
	}
	err := order.Find()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n", order)
}

func TestOrder_Save(t *testing.T) {
	order := Order{
		ID:     1,
		Status: OrderStatusPayed,
		PayAt:  time.Now(),
	}
	err := order.Save()
	if err != nil {
		t.Error(err)
	}
}

func TestOrderList(t *testing.T) {
	order := Order{
		ID:       1,
		UserID:   "test",
		FileList: []int64{1, 2, 3},
		Status:   OrderStatusWaitPay,
		TotalFee: 30,
	}
	err := order.Create()
	if err != nil {
		t.Error(err)
	}
	order.ID = 2
	err = order.Create()
	if err != nil {
		t.Error(err)
	}

	list, err := OrderList("test")
	if err != nil {
		t.Error(err)
	}

	// 输出list
	for _, v := range list {
		fmt.Printf("%+v\n", v)
	}

	// 删除order
	err = order.Delete()
	if err != nil {
		t.Error(err)
	}
	order.ID = 1
	err = order.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestOrder_Delete(t *testing.T) {
	order := Order{
		ID: 1,
	}
	err := order.Delete()
	if err != nil {
		t.Error(err)
	}
}
