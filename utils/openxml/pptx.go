package openxml

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
)

type PptxPropsAppXML struct {
	Template             string  `xml:"Template"`
	TotalTime            int     `xml:"TotalTime"` // 分钟
	Words                int     `xml:"Words"`
	Application          string  `xml:"Application"`
	PresentationFormat   string  `xml:"PresentationFormat"`
	Paragraphs           int     `xml:"Paragraphs"`
	Slides               int     `xml:"Slides"` // 张数
	Notes                int     `xml:"Notes"`
	HiddenSlides         int     `xml:"HiddenSlides"`
	MMClips              int     `xml:"MMClips"`
	ScaleCrop            string  `xml:"ScaleCrop"`
	Company              string  `xml:"Company"`
	LinksUpToDate        bool    `xml:"LinksUpToDate"`
	CharactersWithSpaces int     `xml:"CharactersWithSpaces"`
	SharedDoc            bool    `xml:"SharedDoc"`
	HyperlinksChanged    bool    `xml:"HyperlinksChanged"`
	AppVersion           float32 `xml:"AppVersion"`
}

func GetPptxPages(filePath string) (int, error) {
	// 打开压缩文件
	zipFile, err := zip.OpenReader(filePath)
	if err != nil {
		return 0, err
	}
	// 解压文件
	for _, f := range zipFile.File {
		if f.Name == "docProps/app.xml" {
			// 读取文件内容
			rc, err := f.Open()
			if err != nil {
				return 0, err
			}
			body, err := io.ReadAll(rc)
			if err != nil {
				return 0, err
			}
			// 解析文件内容
			docPropsAppXML := PptxPropsAppXML{}
			if err = xml.Unmarshal(body, &docPropsAppXML); err != nil {
				return 0, err
			}
			if err := rc.Close(); err != nil {
				return 0, err
			}
			return docPropsAppXML.Slides, nil
		}
	}
	return 0, errors.New("docProps/app.xml not found")
}
