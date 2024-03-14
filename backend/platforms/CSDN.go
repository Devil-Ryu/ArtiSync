package platforms

import (
	"ArtiSync/backend/api"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// RodCSDN CSDN机器人
type RodCSDN struct {
	Article   *api.Article // 待上传文章
	CheckTime int          // 检查时间
}

// NewRodCSDN 初始化
func NewRodCSDN() *RodCSDN {
	return &RodCSDN{CheckTime: 2}
}

// SetArticle 设置机器人要上传的文章
func (csdn *RodCSDN) SetArticle(Article *api.Article) {
	csdn.Article = Article
}

// Login 登录
func (csdn *RodCSDN) Login(browser *rod.Browser) []*proto.NetworkCookie {
	fmt.Println("进入登录页面: https://passport.csdn.net/login?code=applets")
	var loginCookies []*proto.NetworkCookie
	// 访问登录页面
	browser.MustPage("https://passport.csdn.net/login?code=applets")

	fmt.Println("MustNavigate")
	// 监听是否登录成功
	for true {
		pages, _ := browser.Pages()
		fmt.Println("当前页面数: ", len(pages), " , 登录状态: 等待登录")
		// 如果页面全部关闭，则推出
		if len(pages) == 0 {
			break
		} else {
			targetPage, err := pages.FindByURL("www.csdn.net")
			if err == nil {
				fmt.Println("当前页面数: ", len(pages), " , 登录状态: 登录成功")
				loginCookies = targetPage.MustCookies() // 获取cookie，终止监听
				break
			}
			time.Sleep(time.Duration(csdn.CheckTime) * time.Second)
		}
	}

	return loginCookies
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
	page.MustWaitStable().MustElement(imageUpload2CSS).MustSetFiles("/Users/ryu/Documents/test/images/Snipaste_2023-09-21_13-53-22.jpg").WaitInvisible()

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
	createArticleCSS := "div.toolbar-btn:nth-child(6) > a:nth-child(1)"
	contentAreaCSS := "body > div.app.app--light > div.layout > div.layout__panel.flex.flex--row > div > div.layout__panel.flex.flex--row > div.layout__panel.layout__panel--editor > div.editor > pre"
	// imageUploadCSS := "body > div.app.app--light > div.layout > div.layout__panel.flex.flex--row > div > div.layout__panel.layout__panel--navigation-bar.clearfix > nav > div.scroll-box > div:nth-child(1) > div:nth-child(14) > button"
	// imageUpload2CSS := "#pane-upimg > div > div > input[type=file]"

	// 设置浏览器
	// 设置检查时间
	csdn.CheckTime = 2

	// 设置浏览器启动器
	browserPath, ok := launcher.LookPath() // 查找浏览器路径
	var launcherT *launcher.Launcher
	fmt.Println("启用本地浏览器: ", ok)
	if ok { // 如果找到本地浏览器，则使用
		launcherT = launcher.New().Bin(browserPath) // 设置浏览器
	} else { // 如果未找到则不使用本地浏览器
		launcherT = launcher.New() // 设置浏览器
	}
	launcherT.Headless(false) // 设置浏览器参数
	defer launcherT.Cleanup() // cleanup

	// 启动浏览器
	browser := rod.New().NoDefaultDevice().ControlURL(launcherT.MustLaunch()).MustConnect() // 打开浏览器
	fmt.Println("浏览器已启动")

	// 设置cookie
	// cookies := csdn.Login(browser)
	// csdn.saveToJSON(cookies)
	cookies, _ := csdn.readCookies()
	browser.SetCookies(cookies)

	/*创建文章*/
	// 导航到首页
	page := browser.MustPage("https://www.csdn.net")
	// 点击发布文章按钮
	page.MustElement(createArticleCSS).MustClick()
	contentAreaEl := page.MustElement(contentAreaCSS)

	// 清除内容
	csdn.clearContent(contentAreaEl)

	/*上传图片*/
	// 遍历图片列表，上传图片
	for index, imageInfo := range csdn.Article.MarkdownTool.ImagesInfo {
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
	fmt.Println(strings.Join(csdn.Article.MarkdownTool.MarkdownLines, "\n"))
	// 关闭浏览器
	// launcherT.Kill()
	fmt.Println("浏览器关闭")

}

// saveToJSon 保存为json
func (csdn *RodCSDN) saveToJSON(cookies []*proto.NetworkCookie) {
	// 先序列化
	jsonStr, err := json.Marshal(cookies)
	if err != nil {
		fmt.Println("SaveJSONFile  marshal platform error: %w", err)
	}

	file, err := os.Create("file.json")
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

func (csdn *RodCSDN) readCookies() (cookies []*proto.NetworkCookieParam, err error) {
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
