package api

import (
	artlog "ArtiSync/backend/logger"
	"ArtiSync/backend/models"
	"ArtiSync/backend/utils"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ATController 文章控制器
type ATController struct {
	ArticleList   []Article     // 文章列表
	ProxyURL      *url.URL      // 网络代理
	RequestSleep  int           // 请求一次休眠时间
	ImageReadType string        // 图片读取方式
	DBController  *DBController // 数据库控制器
	artlog.Logger
}

// Article 文章结构体
type Article struct {
	Title         string                          // 文章名称
	Status        string                          // 文章状态（等待中，上传中，完成）
	MarkdownTool  utils.MarkdownTool              // Markdown分析工具，分析的文章在这里
	Progress      float32                         // 运行进度
	BasicInfo     ArticleBasicInfo                // 文章基本信息
	PlatformsInfo map[string]*ArticlePlatformInfo // 文章平台信息
}

// ArticleBasicInfo 文章基本信息
type ArticleBasicInfo struct {
	Title             string   // 文章名称
	Type              string   // 文章类型
	ImageNum          int      // 图片数量
	InterfaceNum      int      // 接口数量
	InterfaceRecordID string   // 接口记录ID
	Platforms         []string // 平台名称
}

// ArticlePlatformInfo 文章平台信息
type ArticlePlatformInfo struct {
	Index                  int     // 平台序号
	InterfaceActualRunNum  int     // 实际接口运行数量
	InterfacePredictRunNum int     // 预测接口运行数量
	RunProgress            float32 // 运行进度
	RunStatus              string  // 运行状态
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

	// 获取配置文件路径
	configFilePath, err := a.DBController.GetConfigFilePath()
	if err != nil {
		a.Print(artlog.ERROR, "文章控制器", fmt.Errorf("获取配置文件失败: %w", err).Error())
	}

	// 加载config
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

	// 设置NetController的requestSleep
	requestSleep, ok := config["requestSleep"].(float64)
	if !ok {
		err = fmt.Errorf("requestSleep解析错误: %s", config["requestSleep"])
		a.Print(artlog.ERROR, "数据控制器", err.Error())

	}

	a.RequestSleep = int(requestSleep)

	// 设置图片读取方式
	fmt.Printf("imageReadType:%s, %T\n", config["imageReadType"], config["imageReadType"])
	imageReadType, ok := config["imageReadType"].(string)
	if !ok {
		err = fmt.Errorf("imageReadType解析错误: %s", config["imageReadType"])
		a.Print(artlog.ERROR, "数据控制器", err.Error())

	}

	fmt.Printf("a.ImageReadType %s\n", a.ImageReadType)

	a.ImageReadType = imageReadType

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
				Title:         file.Name()[:len(file.Name())-3],
				Status:        "等待中",
				PlatformsInfo: make(map[string]*ArticlePlatformInfo),
			}
			// article.MarkdownTool.Startup(a.Ctx) // 设置context，以启用runtime

			// 设置文章路径以及图片路径以及图片读取方式
			article.MarkdownTool.MarkdownPath = path.Join(filePath, file.Name())
			article.MarkdownTool.ImagePath = imagePath
			article.MarkdownTool.ImageReadType = a.ImageReadType

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
			// a.GenArticleDetail(len(a.ArticleList) - 1) // 加载后生成，该文章的基本信息
		}
	}

	return utils.ResponseJSON{StatusCode: 200, Data: a.ArticleList, Message: ""}
}

// GenArticleDetail 获取文章基本的信息，用作展示
func (a *ATController) GenArticleDetail(articleIndex int) {
	article := &a.ArticleList[articleIndex]

	article.BasicInfo = ArticleBasicInfo{}                        // 清空基本信息
	article.PlatformsInfo = make(map[string]*ArticlePlatformInfo) // 清空基本信息

	// 基本信息
	article.BasicInfo.Title = article.Title
	article.BasicInfo.ImageNum = len(article.MarkdownTool.ImagesInfo)
	article.BasicInfo.Type = "MarkDown"

	// 平台信息
	platforms, _ := a.DBController.GetPlatforms(map[string]interface{}{"Disabled": false})

	for index, platform := range platforms {
		// 基础信息中添加平台名称
		article.BasicInfo.Platforms = append(article.BasicInfo.Platforms, platform.Name)

		// 收集平台信息
		// platformInfo := &ArticlePlatformInfo{}
		article.PlatformsInfo[platform.Name] = &ArticlePlatformInfo{}
		article.PlatformsInfo[platform.Name].Index = index
		article.PlatformsInfo[platform.Name].RunStatus = utils.PublishWating
		for _, interfacesInfo := range platform.Interfaces {
			// 不统计禁用的接口
			if interfacesInfo.Disabled {
				continue
			}
			// 如果为组，则遍历子接口
			if interfacesInfo.IsGroup {
				// 计算接口组需要运行的次数
				// 如果是不读取图片，即图片上传接口不执行

				if interfacesInfo.Type == "images" {
					if a.ImageReadType != "NONE" { // 如果执行图片读取接口，则算入
						article.PlatformsInfo[platform.Name].InterfacePredictRunNum += len(interfacesInfo.Children) * len(article.MarkdownTool.ImagesInfo)
					}
				} else {
					article.PlatformsInfo[platform.Name].InterfacePredictRunNum = len(interfacesInfo.Children)
				}
			} else { // 如果不为组，则加入列表
				// 计算接口需要运行的次数
				if interfacesInfo.Type == "images" {
					if a.ImageReadType != "NONE" { // 如果执行图片读取接口，则算入
						article.PlatformsInfo[platform.Name].InterfacePredictRunNum += len(article.MarkdownTool.ImagesInfo)
					}
				} else {
					article.PlatformsInfo[platform.Name].InterfacePredictRunNum++
				}
			}

		}
		article.BasicInfo.InterfaceNum += article.PlatformsInfo[platform.Name].InterfacePredictRunNum

		// article.PlatformsInfo[platform.Name] = platformInfo
	}
}

// UpdateArticleDetail 更新文章详情
func (a *ATController) UpdateArticleDetail(articleIndex int, platform models.Platform) {
	if !platform.Disabled {
		platformInfo := a.ArticleList[articleIndex].PlatformsInfo[platform.Name]
		a.ArticleList[articleIndex].PlatformsInfo[platform.Name].RunProgress = float32(platformInfo.InterfaceActualRunNum+1) / float32(platformInfo.InterfacePredictRunNum+1) * 100 // 分子分母同时加一，防止分母为0的情况

		// 更新
		a.ArticleList[articleIndex].Progress = 0
		for platformName := range a.ArticleList[articleIndex].PlatformsInfo {
			a.ArticleList[articleIndex].Progress += float32(a.ArticleList[articleIndex].PlatformsInfo[platformName].InterfaceActualRunNum)
		}
		a.ArticleList[articleIndex].Progress = a.ArticleList[articleIndex].Progress / float32(a.ArticleList[articleIndex].BasicInfo.InterfaceNum) * 100
	}

	runtime.EventsEmit(a.Ctx, "UpdateArticleDetail", articleIndex, a.ArticleList[articleIndex].BasicInfo, a.ArticleList[articleIndex].PlatformsInfo, a.ArticleList[articleIndex].Status, a.ArticleList[articleIndex].Progress)
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

// // Run 运行
// func (a *ATController) Run() utils.ResponseJSON {
// 	// 每次运行生成一个记录ID
// 	a.InterfaceRecordID = models.GenerateRandomKey(8)
// 	// 查询启用的平台，返回平台列表TODO（后续可修改为查询文章对应平台）
// 	platforms, err := a.DBController.GetPlatforms(map[string]interface{}{"Disabled": false})
// 	if err != nil {
// 		errMsg := fmt.Errorf("平台获取错误: %w", err).Error()
// 		a.Print(artlog.ERROR, "文章控制器", errMsg)
// 		return utils.ResponseJSON{StatusCode: 500, Data: nil, Message: errMsg}
// 	}

// 	// 遍历文章列表
// 	for articleIndex, aritcle := range a.ArticleList {
// 		a.GenArticleDetail(articleIndex)                      // 运行的时候，重置基本信息, 每个文章只调用一次该函数
// 		a.ArticleList[articleIndex].Status = utils.Publishing // 更新文章状态

// 		// 遍历平台列表
// 		platformUploadSuccessCount := 0 // 平台成功上传计数
// 		for _, platform := range platforms {
// 			// 获取文章对应平台信息
// 			a.ArticleList[articleIndex].PlatformsInfo[platform.Name].RunStatus = utils.Running // 更新平台状态
// 			a.UpdateArticleDetail(articleIndex, platform)                                      // 更新平台状态
// 			err = a.ArticleList[articleIndex].MarkdownTool.AnalyzeMarkdown()                   // 每个平台，重新分析markdown数据（重置MaikdownLines状态, imagesInfo状态）
// 			if err != nil {
// 				errMsg := fmt.Errorf("Markdown文章分析错误[%s]: %w", aritcle.Title, err).Error()
// 				a.Print(artlog.ERROR, "文章控制器", errMsg)
// 				a.ArticleList[articleIndex].Status = utils.PublishedFailed // 更新文章状态
// 				a.UpdateArticleDetail(articleIndex, platform)              // 更新文章状态
// 				continue                                                   // 如果报错则开始下一个平台
// 			}
// 			// 每一个平台单独一个网络控制器
// 			a.CurNetController = utils.NewNetWorkController()
// 			a.CurNetController.SetProxyURL(a.ProxyURL)    // 设置代理
// 			a.CurNetController.SleepTime = a.RequestSleep // 设置休眠时间
// 			// 将文章转化为字符串放入网络请求池中
// 			a.CurNetController.ResponsePool["ART-TITLE"] = a.ArticleList[articleIndex].Title
// 			a.CurNetController.ResponsePoolType["ART-TITLE"] = "TEXT"
// 			a.CurNetController.ResponsePool["ART-CONTENT-STR"] = strings.Join(a.ArticleList[articleIndex].MarkdownTool.MarkdownLines, "\n")
// 			a.CurNetController.ResponsePoolType["ART-CONTENT-STR"] = "TEXT"

// 			a.Print(artlog.DEBUG, "网络控制器", fmt.Sprintf(`{"platformName":"%s", "articleName":"%s", "content":%s}`, platform.Name, aritcle.Title, strings.Join(a.ArticleList[articleIndex].MarkdownTool.MarkdownLines, "\n")))
// 			// 遍历接口
// 			var interfaceErr error // 查看遍历接口时是否有接口出错
// 			for _, interfaceInfo := range platform.Interfaces {
// 				// 如果接口被禁用则不执行下述操作
// 				if interfaceInfo.Disabled {
// 					continue
// 				}
// 				// 判断接口是否为组（组执行： 判断接口是否为图片->循环执行children， 否->执行children）
// 				if interfaceInfo.IsGroup {
// 					// 如果为组则执行子接口列表
// 					interfaceErr = a.runInterfaces(a.CurNetController, platform, articleIndex, interfaceInfo.Children, interfaceInfo.Type)
// 					if interfaceErr != nil {
// 						errMsg := fmt.Errorf("接口运行错误[%s]: %w", interfaceInfo.Name, err).Error()
// 						a.Print(artlog.ERROR, "文章控制器", errMsg)
// 					}
// 				} else {
// 					// 如果不为组则执行接口（接口执行： 判断接口是否为图片->循环执行接口， 否->执行接口）
// 					interfaceErr = a.runInterfaces(a.CurNetController, platform, articleIndex, []models.Interface{interfaceInfo}, interfaceInfo.Type)
// 					if interfaceErr != nil {
// 						errMsg := fmt.Errorf("接口运行错误[%s]: %w", interfaceInfo.Name, err).Error()
// 						a.Print(artlog.ERROR, "文章控制器", errMsg)
// 					}
// 				}
// 			}
// 			if interfaceErr != nil { // 若遍历接口时有报错，则平台运行失败
// 				a.ArticleList[articleIndex].PlatformsInfo[platform.Name].RunStatus = utils.PublishedFailed // 更新平台状态
// 			} else { // 若不报错则上传成功，上传成功的平台计数+1
// 				a.ArticleList[articleIndex].PlatformsInfo[platform.Name].RunStatus = utils.PublishedSuccess // 更新平台状态
// 				platformUploadSuccessCount++
// 			}
// 			a.UpdateArticleDetail(articleIndex, platform) // 更新平台状态
// 		}

// 		// 遍历完所有平台后更新文章状态
// 		if platformUploadSuccessCount == len(platforms) {
// 			a.ArticleList[articleIndex].Status = utils.PublishedSuccess
// 		} else {
// 			a.ArticleList[articleIndex].Status = utils.PublishedFailed
// 		}

// 		a.UpdateArticleDetail(articleIndex, models.Platform{Disabled: true}) // 更新文章详情

// 	}
// 	// 运行结束
// 	return utils.ResponseJSON{StatusCode: 200, Data: nil, Message: "运行结束"}
// }
