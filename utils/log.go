package utils

import (
	"io"
	"os"
)

func HttpClose(hp io.ReadCloser) {
	err := hp.Close()
	if err != nil {

	}
}

func FileClose(fp *os.File) {
	err := fp.Close()
	if err != nil {

	}
}
