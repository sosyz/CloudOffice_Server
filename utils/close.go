package utils

import (
	"io"
	"os"
	"sonui.cn/cloudprint/utils/log"
)

func HttpClose(hp io.ReadCloser) {
	err := hp.Close()
	if err != nil {
		log.NewError(2, err.Error())
	}
}

func FileClose(fp *os.File) {
	err := fp.Close()
	if err != nil {
		log.NewError(2, err.Error())
	}
}
