package openxml

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
)

type DocxPropsAppXML struct {
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

// GetDocxInfo 获取docx文件信息
func GetDocxInfo(filePath string) (*DocxPropsAppXML, error) {
	// 打开压缩文件
	zipFile, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	// 解压文件
	for _, f := range zipFile.File {
		if f.Name == "docProps/app.xml" {
			// 读取文件内容
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			body, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}
			// 解析文件内容
			docxPropsAppXML := DocxPropsAppXML{}
			if err = xml.Unmarshal(body, &docxPropsAppXML); err != nil {
				return nil, err
			}
			if err := rc.Close(); err != nil {
				return nil, err
			}
			return &docxPropsAppXML, nil
		}
	}
	return nil, errors.New("docProps/app.xml not found")
}
