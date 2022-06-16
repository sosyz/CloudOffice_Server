package utils

import "io"

func HttpClose(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {

	}
}
