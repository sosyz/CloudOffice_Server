package openxml

import "testing"

func TestGetDocxPages(t *testing.T) {
	pages, err := GetDocxPages("D:\\test.docx")
	if err != nil {
		t.Error(err)
	}
	t.Log(pages)
}
