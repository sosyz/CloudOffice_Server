package models

import (
	"fmt"
	"testing"
	"time"
)

func TestCache_All(t *testing.T) {
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	v.Name = "test"
	v.Age = 1
	cache := Cache{
		Key:      "test",
		Value:    v,
		ExpireAt: time.Now().Add(time.Hour * 24),
	}
	err := cache.Create()
	if err != nil {
		t.Error(err)
	}
	err = cache.Find()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", cache.Value)
	err = cache.Delete()
	if err != nil {
		t.Error(err)
	}
}
