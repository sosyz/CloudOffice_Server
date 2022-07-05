package openxml

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
)

type DocPropsAppXML struct {
	Template             string  `xml:"Template"`
	TotalTime            int     `xml:"TotalTime"`
	Pages                int     `xml:"Pages"`
	Words                int     `xml:"Words"`
	Characters           int     `xml:"Characters"`
	Application          string  `xml:"Application"`
	DocSecurity          int     `xml:"DocSecurity"`
	Lines                int     `xml:"Lines"`
	Paragraphs           int     `xml:"Paragraphs"`
	ScaleCrop            bool    `xml:"ScaleCrop"`
	Company              string  `xml:"Company"`
	LinksUpToDate        bool    `xml:"LinksUpToDate"`
	CharactersWithSpaces int     `xml:"CharactersWithSpaces"`
	SharedDoc            bool    `xml:"SharedDoc"`
	HyperlinksChanged    bool    `xml:"HyperlinksChanged"`
	AppVersion           float32 `xml:"AppVersion"`
}

// GetDocxPages 获取文件页数
func GetDocxPages(filePath string) (int, error) {
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
			docPropsAppXML := DocPropsAppXML{}
			if err = xml.Unmarshal(body, &docPropsAppXML); err != nil {
				return 0, err
			}
			if err := rc.Close(); err != nil {
				return 0, err
			}
			return docPropsAppXML.Pages, nil
		}
	}
	return 0, errors.New("docProps/app.xml not found")
}
