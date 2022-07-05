package openxml

import "testing"

func TestGetPptxPages(t *testing.T) {
	pages, err := GetPptxPages("C:\\Users\\Sonui\\OneDrive\\比赛\\云打印\\CloudOffice.pptx")
	if err != nil {
		t.Error(err)
	}
	t.Log(pages)
}
