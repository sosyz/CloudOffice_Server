package models

import (
	"fmt"
	"testing"
)

func TestFile_Create(t *testing.T) {
	file := File{
		Fid:  1,
		Name: "test",
		Path: "test",
		Size: 1,
	}
	fmt.Printf("%+v\n", file)
	err := file.Create()
	if err != nil {
		t.Error(err)
	}
}

func TestFile_Exist(t *testing.T) {
	file := File{
		Fid: 1,
	}
	ok, err := file.Exist()
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Error("not exist")
	}
}

func TestFile_Find(t *testing.T) {
	file := File{
		Fid: 1,
	}
	err := file.Find()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n", file)
}

func TestFile_Save(t *testing.T) {
	file := File{
		Fid:  1,
		Name: "test2",
	}
	err := file.Save()
	if err != nil {
		t.Error(err)
	}
}

func TestFile_Delete(t *testing.T) {
	file := File{
		Fid: 1,
	}
	err := file.Delete()
	if err != nil {
		t.Error(err)
	}
}
