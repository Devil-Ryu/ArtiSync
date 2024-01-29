package api

import (
	artlog "ArtiSync/backend/logger"
	"ArtiSync/backend/models"
	"ArtiSync/backend/utils"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ATController 文章控制器
type ATController struct {
	ArticleList  []Article     // 文章列表
	ProxyURL     *url.URL      // 网络代理
	DBController *DBController // 数据库控制器
	artlog.Logger
}

// Article 文章结构体
type Article struct {
	Title            string                              // 文章名称
	Status           string                              // 文章状态（等待中，上传中，完成）
	MarkdownTool     utils.MarkdownTool                  // Markdown分析工具，分析的文章在这里
	ImageProgressMap map[string]map[string]ImageProgress // 平台->图片->状态
}

// ImageProgress 图片进度
type ImageProgress struct {
	PlatformName string // 平台名称
	ImageName    string // 图片名称
	UploadURL    string // 上传URL
	Status       string // 上传状态（待上传，成功，失败，上传中）
	Message      string // 图片错误信息
}

// NewATController 初始化
func NewATController() *ATController {
	return &ATController{}
}

// SetDBController 设置DBController
func (a *ATController) SetDBController(dbController *DBController) {
	a.DBController = dbController
}

// InitConfig 初始化配置
func (a *ATController) InitConfig() (err error) {

	configFilePath, err := a.DBController.GetConfigFilePath()
	if err != nil {
		a.Print(artlog.ERROR, "文章控制器", fmt.Errorf("获取配置文件失败: %w", err).Error())
	}

	config, err := a.DBController.LoadJSONFile(configFilePath)
	if err != nil {
		a.Print(artlog.ERROR, "文章控制器", fmt.Errorf("初始化配置错误: %w", err).Error())
	}

	// 设置DBController数据库地址
	dbPath, ok := config["dbPath"].(string)
	if !ok {
		a.Print(artlog.ERROR, "数据库地址解析错误", dbPath)
		err = fmt.Errorf("数据库地址解析错误: %s", dbPath)
	}
	a.DBController.Connect(dbPath)
	a.Print(artlog.DEBUG, "数据控制器", fmt.Sprintf("设置数据库地址为: %s", dbPath))

	// 设置ATController的proxyURL
	proxyURL, ok := config["proxyURL"].(string)
	if !ok {
		err = fmt.Errorf("proxyURL解析错误: %s", proxyURL)
		a.Print(artlog.ERROR, "数据控制器", err.Error())

	}
	a.SetProxyURL(proxyURL)
	return err

}

// OpenDir 打开文件夹
func (a *ATController) OpenDir(defaultPath string) utils.ResponseJSON {
	opts := runtime.OpenDialogOptions{
		DefaultDirectory: filepath.Dir(defaultPath),
	}
	selected, err := runtime.OpenDirectoryDialog(a.Ctx, opts)
	if err != nil {
		errMsg := fmt.Errorf("打开文件夹失败: %w", err).Error()
		a.Print(artlog.ERROR, "文章控制器", errMsg)
		return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
	}
	// a.Print(artlog.DEBUG, "文章控制器", "选择文件夹: "+selected)
	return utils.ResponseJSON{StatusCode: 200, Data: selected, Message: ""}
}

// OpenFile 打开文件
func (a *ATController) OpenFile(defaultPath string) utils.ResponseJSON {

	opts := runtime.OpenDialogOptions{
		DefaultDirectory: filepath.Dir(defaultPath),
	}
	selected, err := runtime.OpenFileDialog(a.Ctx, opts)
	if err != nil {
		errMsg := fmt.Errorf("打开文件失败: %w", err).Error()
		a.Print(artlog.ERROR, "文章控制器", errMsg)
		return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
	}
	// a.Print(artlog.DEBUG, "文章控制器", "选择文件: "+selected)
	return utils.ResponseJSON{StatusCode: 200, Data: selected, Message: ""}
}

// SaveFile 保存文件
func (a *ATController) SaveFile(defaultPath string, defaultFileName string, title string) utils.ResponseJSON {
	opts := runtime.SaveDialogOptions{
		DefaultDirectory:     filepath.Dir(defaultPath),
		DefaultFilename:      defaultFileName,
		CanCreateDirectories: true,
		Title:                title,
	}
	selected, err := runtime.SaveFileDialog(a.Ctx, opts)
	if err != nil {
		errMsg := fmt.Errorf("保存文件失败: %w", err).Error()
		a.Print(artlog.ERROR, "文章控制器", errMsg)
		return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
	}
	// a.Print(artlog.DEBUG, "文章控制器", "选择文件夹: "+selected)
	return utils.ResponseJSON{StatusCode: 200, Data: selected, Message: ""}
}

// GetLogHistory 返回日志历史信息
func (a *ATController) GetLogHistory() utils.ResponseJSON {
	return utils.ResponseJSON{StatusCode: 200, Data: a.History, Message: ""}
}

// LoadArticles 加载文章
func (a *ATController) LoadArticles(filePath string, imagePath string) utils.ResponseJSON {

	// 加载前清除数据
	a.ArticleList = []Article{}

	files, err := os.ReadDir(filePath)
	if err != nil {
		errMsg := fmt.Errorf("读取文章文件夹失败: %w", err).Error()
		a.Print(artlog.ERROR, "文章控制器", errMsg)
		return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
	}

	// 遍历文件夹，获取文章信息
	for _, file := range files {
		if path.Ext(file.Name()) == ".md" {
			article := Article{
				Title:            file.Name()[:len(file.Name())-3],
				Status:           "等待中",
				ImageProgressMap: make(map[string]map[string]ImageProgress),
			}
			// article.MarkdownTool.Startup(a.Ctx) // 设置context，以启用runtime

			// 设置文章路径以及图片路径
			article.MarkdownTool.MarkdownPath = path.Join(filePath, file.Name())
			article.MarkdownTool.ImagePath = imagePath

			// 开始分析文章
			err = article.MarkdownTool.AnalyzeMarkdown()
			if err != nil {
				errMsg := fmt.Errorf("文章分析失败[%s]: %w", article.Title, err).Error()
				a.Print(artlog.ERROR, "文章控制器", errMsg)
				return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
			}
			article.MarkdownTool.Ctx = nil // Ctx不返回前端
			a.ArticleList = append(a.ArticleList, article)

			// 信息传回前端
			a.Print(artlog.DEBUG, "文章控制器", fmt.Sprintf(`{"fileName":"%s", "filePath": "%s", "imagePath":"%s", "imageNum":"%d"}`, file.Name(), filePath, imagePath, len(article.MarkdownTool.ImagesInfo)))

		}
	}

	// articleList := []string{}
	// for index, article := range a.ArticleList {
	// 	articleList = append(articleList, fmt.Sprintf("%d: %s[%03d]", index, article.Title, len(article.MarkdownTool.ImagesInfo)))
	// }
	// a.Print(artlog.INFO, "文章控制器", fmt.Sprintf("加载文章: %s", strings.Join(articleList, ", ")))
	return utils.ResponseJSON{StatusCode: 200, Data: a.ArticleList, Message: ""}
}

// GetArticleImageProgress 获取文章的图片信息列表
func (a *ATController) GetArticleImageProgress(articleIndex int) utils.ResponseJSON {
	// 查询启用的平台，返回平台列表
	// 获取文章
	article := a.ArticleList[articleIndex]
	// 获取平台

	platforms, err := a.DBController.GetPlatforms(map[string]interface{}{"Disabled": false})
	if err != nil {
		return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: fmt.Sprintf("GetArticleImageProgress error: %s", err.Error())}
	}

	// 根据启用的平台，和加载的文章生成一个文章的图片进度列表,并根据文章的ImageProgressMap，更新文章状态
	imageProgressList := []ImageProgress{}
	for _, imageInfo := range article.MarkdownTool.ImagesInfo {
		for _, platform := range platforms {
			satus := article.ImageProgressMap[platform.Name][imageInfo.URL].Status
			if article.ImageProgressMap[platform.Name][imageInfo.URL].Status == "" {
				satus = "等待中"
			}
			imageProgressList = append(imageProgressList, ImageProgress{
				PlatformName: platform.Name,
				ImageName:    imageInfo.URL,
				UploadURL:    article.ImageProgressMap[platform.Name][imageInfo.URL].UploadURL,
				Status:       satus,
				Message:      article.ImageProgressMap[platform.Name][imageInfo.URL].Message,
			})
		}

	}

	return utils.ResponseJSON{StatusCode: 200, Data: imageProgressList, Message: ""}

}

// UpdateArticleStatus 更新文章状态
// func (a *ATController) UpdateArticleStatus(articleIndex int, imageProgress ImageProgress) {
func (a *ATController) UpdateArticleStatus(articleIndex int, data interface{}) {
	// 更新文章Map中图片的状态
	imageProgress, ok := data.(ImageProgress)
	if ok {
		// 获取子Map
		subMap := a.ArticleList[articleIndex].ImageProgressMap[imageProgress.PlatformName]
		if len(subMap) == 0 {
			subMap = map[string]ImageProgress{}
		}
		// 更新子Map
		subMap[imageProgress.ImageName] = imageProgress
		a.ArticleList[articleIndex].ImageProgressMap[imageProgress.PlatformName] = subMap
	}
	// 更新完后通知前端
	response := a.GetArticleImageProgress(articleIndex)
	if response.StatusCode != 200 {
		a.Print(artlog.ERROR, "文章控制器", response.Message)
	}
	runtime.EventsEmit(a.Ctx, "UpdateArticleStatus", articleIndex, response.Data, a.ArticleList[articleIndex].Status)
}

func (a *ATController) runInterfaces(netController *utils.NetWorkController, platform models.Platform, articleIndex int, interfaces []models.Interface, interfaceType string) (err error) {
	article := a.ArticleList[articleIndex]
	if interfaceType == "images" {
		// 遍历图片
		for index, imageInfo := range article.MarkdownTool.ImagesInfo {
			// 遍历每个接口(穿过来的一般是children)
			for _, interfaceInfo := range interfaces {
				// 接口放入文件参数
				tmpInterfaceInfo := interfaceInfo
				if tmpInterfaceInfo.Type == "images" { // 如果子接口为图片接口
					tmpInterfaceInfo.RequestBody = append(interfaceInfo.RequestBody, models.Body{IsFile: true, FileName: imageInfo.URL, Key: "file", Value: string(imageInfo.Image)})
				}
				netController.SetInterface(tmpInterfaceInfo)

				response := netController.Run()
				if response.StatusCode == 200 {
					responseStatus := "上传成功"
					a.Print(artlog.DEBUG, "网络控制器", fmt.Sprintf(`{"platformName":"%s", "interfaceName":"%s", "interfaceURL":"%s", "response":%s}`,
						platform.Name,
						tmpInterfaceInfo.Name, tmpInterfaceInfo.RequestURL,
						strconv.Quote(netController.ResponsePool[tmpInterfaceInfo.Serial])))

					imgURL, err := netController.GetResponseMappedValue() // 获取图片链接
					if err != nil {
						a.Print(artlog.ERROR, "网络控制器", fmt.Sprintf(`{"platformName":"%s", "interfaceName":"%s", "interfaceURL":"%s", "response":%s}`,
							platform.Name,
							tmpInterfaceInfo.Name, tmpInterfaceInfo.RequestURL,
							err.Error()))
					}

					article.MarkdownTool.ImagesInfo[index].UploadURL = imgURL // 更新图片链接

					// 此处需要判断上传状态
					imageProgress := ImageProgress{
						PlatformName: platform.Name,
						ImageName:    article.MarkdownTool.ImagesInfo[index].URL,       // 图片名称
						UploadURL:    article.MarkdownTool.ImagesInfo[index].UploadURL, // 图片上传URL
						Status:       responseStatus,
						Message:      response.Message,
					}
					// 更新图片状态
					a.UpdateArticleStatus(articleIndex, imageProgress)

				} else {
					responseStatus := "上传失败"
					a.Print(artlog.ERROR, "网络控制器", fmt.Sprintf(`{"platformName":"%s", "interfaceName":"%s", "interfaceURL":"%s", "response":%s}`,
						platform.Name,
						tmpInterfaceInfo.Name, tmpInterfaceInfo.RequestURL,
						strconv.Quote(response.Message)))

					imageProgress := ImageProgress{
						PlatformName: platform.Name,
						ImageName:    article.MarkdownTool.ImagesInfo[index].URL, // 图片名称
						UploadURL:    response.Message,                           // 图片上传URL
						Status:       responseStatus,
						Message:      response.Message,
					}
					// 更新图片状态
					a.UpdateArticleStatus(articleIndex, imageProgress)
					err = fmt.Errorf(response.Message)
					// return fmt.Errorf(response.Message)
				}
			}

		}
		article.MarkdownTool.ReplaceImages() // 替换链接

		// 将替换后的文章转化为字符串放入网络请求池中
		content := map[string]string{
			"article": strings.Join(article.MarkdownTool.MarkdownLines, "\n"),
			"title":   article.Title,
		}
		data, _ := json.Marshal(content)
		netController.ResponsePool["0"] = string(data)
		netController.ResponsePoolType["0"] = "JSON" // 设置接口返回类型为JSON
	} else {
		// 遍历每个接口并运行（一般不是images，只传一个接口过来）
		for _, interfaceInfo := range interfaces {
			netController.SetInterface(interfaceInfo)
			// netController.Run()
			response := netController.Run()
			if response.StatusCode == 200 {
				a.Print(artlog.DEBUG, "网络控制器", fmt.Sprintf(`{"platformName":"%s", "interfaceName":"%s", "interfaceURL":"%s", "response":%s}`,
					platform.Name,
					interfaceInfo.Name, interfaceInfo.RequestURL,
					strconv.Quote(netController.ResponsePool[interfaceInfo.Serial])))
			} else {
				a.Print(artlog.ERROR, "网络控制器", fmt.Sprintf(`{"platformName":"%s", "interfaceName":"%s", "interfaceURL":"%s", "response":%s}`,
					platform.Name,
					interfaceInfo.Name, interfaceInfo.RequestURL,
					strconv.Quote((response.Message))))
				err = fmt.Errorf(response.Message)
				// return fmt.Errorf(response.Message)
			}
		}
	}

	// 最终返回err
	return err
}

// SetProxyURL 设置代理
func (a *ATController) SetProxyURL(proxURL string) {
	tmp, err := url.Parse(proxURL)
	a.ProxyURL = tmp
	// 空代理设置为nil
	if proxURL == "" {
		a.ProxyURL = nil
	}
	a.Print(artlog.DEBUG, "网络控制器", fmt.Sprintf("设置代理为: %s", a.ProxyURL))
	if err != nil {
		a.Print(artlog.WARN, "网络控制器", fmt.Sprintf("设置代理: %s", err.Error()))
	}

}

// Run 运行
func (a *ATController) Run() utils.ResponseJSON {
	// 查询启用的平台，返回平台列表
	platforms, err := a.DBController.GetPlatforms(map[string]interface{}{"Disabled": false})
	if err != nil {
		errMsg := fmt.Errorf("平台获取错误: %w", err).Error()
		a.Print(artlog.ERROR, "文章控制器", errMsg)
		return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
	}

	// 遍历文章列表
	for articleIndex, aritcle := range a.ArticleList {
		// a.Print(artlog.INFO, "文章控制器", fmt.Sprintf("任务开始: %s", aritcle.Title))
		// 遍历平台列表
		a.ArticleList[articleIndex].Status = "处理中"
		a.UpdateArticleStatus(articleIndex, a.ArticleList[articleIndex].Status)
		for _, platform := range platforms {
			err = a.ArticleList[articleIndex].MarkdownTool.AnalyzeMarkdown() // 每个平台，重新分析markdown数据（重置MaikdownLines状态, imagesInfo状态）
			if err != nil {
				errMsg := fmt.Errorf("Markdown文章分析错误[%s]: %w", aritcle.Title, err).Error()
				a.Print(artlog.ERROR, "文章控制器", errMsg)
				a.ArticleList[articleIndex].Status = "上传失败"
				a.UpdateArticleStatus(articleIndex, a.ArticleList[articleIndex].Status)
				continue // 如果报错则开始下一个平台
				// return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
			}
			// 每一个平台单独一个网络控制器
			netController := utils.NewNetWorkController()
			netController.SetProxyURL(a.ProxyURL) // 设置代理
			// 将文章转化为字符串放入网络请求池中
			content := map[string]string{
				"article": strings.Join(a.ArticleList[articleIndex].MarkdownTool.MarkdownLines, "\n"),
				"title":   a.ArticleList[articleIndex].Title,
			}
			data, err := json.Marshal(content)
			if err != nil {
				errMsg := fmt.Errorf("文章信息转失败[%s]: %w", aritcle.Title, err).Error()
				a.Print(artlog.ERROR, "文章控制器", errMsg)
				a.ArticleList[articleIndex].Status = "上传失败"
				a.UpdateArticleStatus(articleIndex, a.ArticleList[articleIndex].Status)
				continue // 如果报错则开始下一个平台
				// return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
			}
			netController.ResponsePool["0"] = string(data)
			netController.ResponsePoolType["0"] = "JSON" // 设置接口返回类型为JSON
			a.Print(artlog.DEBUG, "网络控制器", fmt.Sprintf(`{"platformName":"%s", "articleName":"%s", "content":%s}`, platform.Name, aritcle.Title, string(data)))

			// 遍历接口
			// 组执行： 判断接口是否为图片->循环执行children， 否->执行children
			// 接口执行： 判断接口是否为图片->循环执行接口， 否->执行接口
			for index, interfaceInfo := range platform.Interfaces {
				// 判断接口是否为组
				if interfaceInfo.IsGroup {
					// 如果为组则执行子接口列表
					err = a.runInterfaces(netController, platform, articleIndex, interfaceInfo.Children, interfaceInfo.Type)
					if err != nil {
						errMsg := fmt.Errorf("接口运行错误[%s]: %w", interfaceInfo.Name, err).Error()
						a.Print(artlog.ERROR, "文章控制器", errMsg)
						a.ArticleList[articleIndex].Status = "上传失败"
						break // 如果报错则跳出
						// return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
					}
				} else {
					// 如果不为组则执行接口
					err = a.runInterfaces(netController, platform, articleIndex, []models.Interface{interfaceInfo}, interfaceInfo.Type)
					if err != nil {
						errMsg := fmt.Errorf("接口运行错误[%s]: %w", interfaceInfo.Name, err).Error()
						a.Print(artlog.ERROR, "文章控制器", errMsg)
						a.ArticleList[articleIndex].Status = "上传失败"
						break // 如果报错则跳出
						// return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
					}
				}
				// 所有接口成功执行后返回上传成功
				if index == (len(platform.Interfaces) - 1) {
					a.ArticleList[articleIndex].Status = "运行完成"
				}

			}
			a.UpdateArticleStatus(articleIndex, a.ArticleList[articleIndex].Status)
			// a.Print(artlog.INFO, "文章控制器", fmt.Sprintf("任务结束: %s", aritcle.Title))
		}

	}
	// 运行结束
	return utils.ResponseJSON{StatusCode: 200, Data: nil, Message: "运行结束"}
}
