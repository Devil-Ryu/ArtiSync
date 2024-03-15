package platforms

import (
	"ArtiSync/backend/api"
	"ArtiSync/backend/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// RodCSDN CSDN机器人
type RodCSDN struct {
	RODController *utils.RODController        // 机器人控制器
	Cookies       []*proto.NetworkCookieParam // cookies
	Article       *api.Article                // 待上传文章
}

// ConfigCSDN 配置文件
type ConfigCSDN struct {
	CookieSavePath           string `json:"CookieSavePath"`           // cookie保存地址
	PlatformName             string `json:"PlatformName"`             // 平台名称
	LoginPageURL             string `json:"LoginPageURL"`             // 登录页
	HomePageURL              string `json:"HomePageURL"`              // 首页
	ProfilePageURL           string `json:"ProfilePageURL"`           // 个人页
	CreateArticleBtnSelector string `json:"CreateArticleBtnSelector"` // 创建文章选择器
	TitleInputSelector       string `json:"TitleInputSelector"`       // 标题输入选择器
	ContentAreaSelector      string `json:"ContentAreaSelector"`      // 内容区域选择器
	ImageUploadStep1Selector string `json:"ImageUploadStep1Selector"` // 图片上传选择器
	ImageUploadStep2Selector string `json:"ImageUploadStep2Selector"` // 图片上传选择器
	SaveArticleBtnSelector   string `json:"SaveArticleBtnSelector"`   // 保存文章选择器
	ProfileNameSelector      string `json:"ProfileNameSelector"`      // 个人简介名称选择器
	ProfileAvatarSelector    string `json:"ProfileAvatarSelector"`    // 个人简介头像选择器

}

// NewRodCSDN 初始化
func NewRodCSDN() *RodCSDN {
	return &RodCSDN{RODController: &utils.RODController{}}
}

// SetArticle 设置机器人要上传的文章
func (csdn *RodCSDN) SetArticle(Article *api.Article) {
	csdn.Article = Article
}

// SetCookies 设置cookies
func (csdn *RodCSDN) SetCookies(Cookies []*proto.NetworkCookieParam) {
	csdn.Cookies = Cookies
}

// LoadConfig 加载配置TODO（将configCSDN内容从本地加载到struct）
func (csdn *RodCSDN) LoadConfig() {

}

// LoginCSDN 登录CSDN后把cookie保存到本地
func (csdn *RodCSDN) LoginCSDN() utils.ResponseJSON {
	// 查看当前是否有浏览器
	if csdn.RODController.Browser == nil {
		csdn.RODController.StartBrowser()
	}

	fmt.Println("进入登录页面: https://passport.csdn.net/login?code=applets")
	var loginCookies []*proto.NetworkCookie
	// 访问登录页面
	csdn.RODController.Browser.MustPage("https://passport.csdn.net/login?code=applets")

	// 监听是否登录成功
	for true {
		pages, _ := csdn.RODController.Browser.Pages()
		fmt.Println("当前页面数: ", len(pages), " , 登录状态: 等待登录")
		// 如果页面全部关闭，则推出
		if len(pages) == 0 {
			fmt.Println("当前页面数: ", len(pages), " , 登录状态: 取消登录")
			break
		} else {
			targetPage, err := pages.FindByURL("www.csdn.net")
			if err == nil {
				fmt.Println("当前页面数: ", len(pages), " , 登录状态: 登录成功")
				loginCookies = targetPage.MustCookies() // 获取cookie，终止监听
				break
			}
			time.Sleep(time.Duration(csdn.RODController.CheckTime) * time.Second)
		}
	}

	// 关闭浏览器
	csdn.RODController.CloseBrowser()
	csdn.SaveToJSON(loginCookies, "CSDN_Cookies.json")

	return utils.ResponseJSON{StatusCode: 200, Data: "CSDN_Cookies.json", Message: ""}
}

// GetProfile 获取信息
func (csdn *RodCSDN) GetProfile() {
	profileURL := "https://i.csdn.net/#/user-center/profile"

	/*设置浏览器*/
	if csdn.RODController.Browser == nil {
		csdn.RODController.StartBrowser() // 启动浏览器
	}

	/*设置浏览器cookies*/
	csdn.RODController.Browser.SetCookies(csdn.Cookies)

	page := csdn.RODController.Browser.MustPage(profileURL)
	name, err := page.MustWaitStable().MustElement(".el-form > ul:nth-child(1) > li:nth-child(1) > div:nth-child(2)").Text()
	if err != nil {
		name = err.Error()
	}
	avater, err := page.MustWaitStable().MustElement(".avatar-hover > img:nth-child(1)").Attribute("src")
	if err != nil {
		name = err.Error()
	}
	fmt.Println("name: ", name)
	fmt.Println("name: ", *avater)

	csdn.RODController.CloseBrowser()

}

// clearContent 清除内容
func (csdn *RodCSDN) clearContent(element *rod.Element) {
	// 找到所有的输入
	subEls := element.MustElementsX(".//div")

	// 清空输入区
	for _, subEl := range subEls {
		subEl.Remove()
	}
}

// UploadImage 上传图片
func (csdn *RodCSDN) UploadImage(page *rod.Page, element *rod.Element, imagePath string) (uploadURL string, err error) {
	// 访问上传图片页面
	imageUploadCSS := "body > div.app.app--light > div.layout > div.layout__panel.flex.flex--row > div > div.layout__panel.layout__panel--navigation-bar.clearfix > nav > div.scroll-box > div:nth-child(1) > div:nth-child(14) > button"
	imageUpload2CSS := "#pane-upimg > div > div > input[type=file]"

	// 上传图片
	page.MustElement(imageUploadCSS).MustClick()
	page.MustWaitStable().MustElement(imageUpload2CSS).MustSetFiles(imagePath).WaitInvisible()

	// 获取图片元素
	imgEl, err := element.MustWaitStable().ElementX("*//img")
	if err != nil {
		return uploadURL, err
	}
	// 获取图片链接
	uploadURL = *imgEl.MustAttribute("src")
	fmt.Println("uploadUrl: ", uploadURL)
	// 获得链接后清除内容
	csdn.clearContent(element)
	return uploadURL, nil
}

// RUN 运行
func (csdn *RodCSDN) RUN() {
	fmt.Println("开始运行: CSDN")
	createArticleCSS := "div.toolbar-btn:nth-child(6) > a:nth-child(1)"
	titleCSS := "body > div.app.app--light > div.layout > div.layout__panel.layout__panel--articletitle-bar > div > div.article-bar__input-box > input"
	contentAreaCSS := "body > div.app.app--light > div.layout > div.layout__panel.flex.flex--row > div > div.layout__panel.flex.flex--row > div.layout__panel.layout__panel--editor > div.editor > pre"
	// imageUploadCSS := "body > div.app.app--light > div.layout > div.layout__panel.flex.flex--row > div > div.layout__panel.layout__panel--navigation-bar.clearfix > nav > div.scroll-box > div:nth-child(1) > div:nth-child(14) > button"
	// imageUpload2CSS := "#pane-upimg > div > div > input[type=file]"
	saveArticleCSS := "body > div.app.app--light > div.layout > div.layout__panel.layout__panel--articletitle-bar > div > div.article-bar__user-box.flex.flex--row > button.btn.btn-save"
	// PublishArticleCSS := "body > div.app.app--light > div.layout > div.layout__panel.layout__panel--articletitle-bar > div > div.article-bar__user-box.flex.flex--row > button.btn.btn-publish"

	/*设置浏览器*/
	if csdn.RODController.Browser == nil {
		csdn.RODController.StartBrowser() // 启动浏览器
	}

	/*设置浏览器cookies*/
	csdn.RODController.Browser.SetCookies(csdn.Cookies)

	/*创建文章*/
	page := csdn.RODController.Browser.MustPage("https://www.csdn.net")          // 导航到首页
	page.MustElement(createArticleCSS).MustClick()                               // 光标移动到发布文章按钮
	contentAreaEl := page.MustElement(contentAreaCSS)                            // 定位内容区域
	csdn.clearContent(contentAreaEl)                                             // 清除现有输入区内容
	page.MustElement(titleCSS).MustSelectAllText().MustInput(csdn.Article.Title) // 输入文章标题

	/*上传图片*/
	for index, imageInfo := range csdn.Article.MarkdownTool.ImagesInfo { // 遍历图片列表，上传图片
		imagePath := path.Join(csdn.Article.MarkdownTool.ImagePath, imageInfo.URL)
		uploadURL, err := csdn.UploadImage(page, contentAreaEl, imagePath)
		if err != nil {
			fmt.Println("err: ", err)
		} else {
			csdn.Article.MarkdownTool.ImagesInfo[index].UploadURL = uploadURL
		}

	}
	/*替换图片*/
	csdn.Article.MarkdownTool.ReplaceImages()
	/*上传文章*/
	for _, line := range csdn.Article.MarkdownTool.MarkdownLines {
		time.Sleep(1 * time.Second)
		for i := 0; i < 20; i++ {
			page.KeyActions().Press(input.ShiftLeft).Press(input.Tab).MustDo()
		}
		contentAreaEl.MustInput(line)
		page.KeyActions().Type(input.Enter).MustDo() // 回车换行
	}

	page.MustElement(saveArticleCSS).MustClick() // 点击保存

	// 关闭浏览器
	csdn.RODController.CloseBrowser()
	fmt.Println("浏览器关闭")

}

// SaveToJSON 保存为json
func (csdn *RodCSDN) SaveToJSON(cookies []*proto.NetworkCookie, savePath string) {
	// 先序列化
	jsonStr, err := json.Marshal(cookies)
	if err != nil {
		fmt.Println("SaveJSONFile  marshal platform error: %w", err)
	}

	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println("SaveJSONFile create file error: %w", err)
	}

	bw := bufio.NewWriter(file)
	_, err = bw.WriteString(string(jsonStr))
	if err != nil {
		fmt.Println("SaveJSONFile write content error: %w", err)
	}
	bw.Flush()
	fmt.Println("Cookie 已经保存到本地")

}

// ReadCookies dddd
func (csdn *RodCSDN) ReadCookies() (cookies []*proto.NetworkCookieParam, err error) {
	dataBytes, err := os.ReadFile("file.json")
	if err != nil {
		return cookies, fmt.Errorf("LoadJSONFile read file error: %w", err)
	}

	// 读取文件内容
	err = json.Unmarshal(dataBytes, &cookies)
	if err != nil {
		return cookies, fmt.Errorf("LoadJSONFile json Unmarshal error: %w", err)
	}
	return cookies, nil
}
