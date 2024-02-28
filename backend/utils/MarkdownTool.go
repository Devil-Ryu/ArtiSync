package utils

import (
	artlog "ArtiSync/backend/logger"
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

// ImgPattern Markdown语法正则表达式
const ImgPattern string = "(\\!\\[(.*?)\\]\\((.*?)\\))"

// ImageInfo 图片信息结构体
type ImageInfo struct {
	Line           int      // 图片在原文的索引
	Title          string   // 图片的名称[title](title)
	URL            string   // 图片的链接[title](url)
	CurURL         string   // 图片的当前链接
	Image          []byte   // 图片二进制流
	UploadURL      string   // 上传后的链接
	UploadPlatform []string // 图片上传的平台（废弃：暂时没用）
}

// MarkdownTool 解析器结构体
type MarkdownTool struct {
	MarkdownLines []string    // 文章数据
	ImagesInfo    []ImageInfo // 图片信息
	MarkdownPath  string      // 文章路径
	ImagePath     string      // 图片路径
	artlog.Logger
}

// AnalyzeMarkdown 分析markdown文章
func (m *MarkdownTool) AnalyzeMarkdown() (err error) {
	// 读取markdown文章
	m.MarkdownLines, err = m.readFile(m.MarkdownPath)
	if err != nil {
		return fmt.Errorf("AnalyzeMarkdown: %w", err)
	}
	// 提取文章中的图片
	m.ImagesInfo, err = m.extractImages(m.MarkdownLines)
	if err != nil {
		return fmt.Errorf("AnalyzeMarkdown: %w", err)
	}
	// 从本地读取图片字节
	err = m.loadImages()
	if err != nil {
		return fmt.Errorf("AnalyzeMarkdown: %w", err)
	}

	return nil
}

// 提取文章中图片信息
func (m *MarkdownTool) extractImages(markdownLines []string) (imagesInfo []ImageInfo, err error) {
	re, err := regexp.Compile(ImgPattern)
	if err != nil {
		return imagesInfo, fmt.Errorf("extractImages: %w", err)
	}
	for index, line := range markdownLines {
		results := re.FindAllStringSubmatch(line, -1)
		for _, result := range results {
			imagesInfo = append(imagesInfo, ImageInfo{Line: index, Title: result[2], URL: result[3], CurURL: result[3]})
		}
	}
	return imagesInfo, nil
}

// 导入本地图片
func (m *MarkdownTool) loadImages() error {
	for index, imgInfo := range m.ImagesInfo {
		imagePath := imgInfo.URL
		image, err := m.ReadImage(path.Join(m.ImagePath, imagePath))
		if err != nil {
			return fmt.Errorf("loadImages: %w", err)
		}
		m.ImagesInfo[index].Image = image
	}
	return nil
}

// ReplaceImages 替换图片
func (m *MarkdownTool) ReplaceImages() {
	for imgIndex, imgInfo := range m.ImagesInfo {
		index := imgInfo.Line
		m.MarkdownLines[index] = strings.Replace(m.MarkdownLines[index], imgInfo.CurURL, imgInfo.UploadURL, 1) // 替换为上传链接
		m.ImagesInfo[imgIndex].CurURL = imgInfo.UploadURL                                                      // 替换后更新当前链接
	}
}

// 读取文件
func (m *MarkdownTool) readFile(fileName string) ([]string, error) {
	var fileList []string
	file, err := os.Open(fileName)

	if err != nil {
		// m.PushLog("error", fmt.Sprintf("Open File Error: %s", err))
		return nil, fmt.Errorf("readFile open: %w", err)
	}

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		fileList = append(fileList, line)
	}

	return fileList, nil
}

// ReadImage 读取图片
func (m *MarkdownTool) ReadImage(imageName string) ([]byte, error) {
	file, err := os.Open(imageName)
	if err != nil {
		// m.PushLog("error", fmt.Sprintf("Open File Error: %s", err))
		return nil, fmt.Errorf("ReadImage open: %w", err)

	}

	// 读取文件内容
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		// m.PushLog("error", fmt.Sprintf("Read Image Error: %s", err))
		return nil, fmt.Errorf("ReadImage read: %w", err)
	}

	return imageBytes, nil

}
