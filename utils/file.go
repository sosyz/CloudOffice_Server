package utils

import (
	"io"
	"os"
)

// Exists 判断文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Read 读取文件内容
func Read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	buf := make([]byte, 1024)
	var content []byte
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		content = append(content, buf[:n]...)
	}
	return content, nil
}

// Write 写入文件内容
func Write(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}

// Append 追加文件内容
func Append(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}

// Remove 删除文件
func Remove(path string) error {
	return os.Remove(path)
}
