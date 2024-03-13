package test

import (
	"ArtiSync/backend/api"
	"ArtiSync/backend/models"
	"ArtiSync/backend/utils"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func TestController(t *testing.T) {
	// 链接数据库
	filePath := "/Users/ryu/Documents/test"
	imagePath := "/Users/ryu/Documents/test/"
	atController := api.NewATController()
	// atController.Startup(ctx)
	atController.LoadArticles(filePath, imagePath)

	a := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(a)
	// atController.Run()

	// a := models.Platform{Name: "test 1111"}

	// response := utils.ResponseJSON{StatusCode: 200, Message: "This is test message"}
	// response.Data = a
	// fmt.Println("response", response)
	// fmt.Println("response", response.Data)
	// data, _ := response.Data.(models.Interface)
	// fmt.Println("data", data)

	// 测试网络模块

}

func TestNetworkController(t *testing.T) {
	controller := api.NewDBController()
	controller.Connect("/Users/ryu/Documents/test.db")
	platform, _ := controller.GetPlatform(models.Platform{ID: 2})
	fmt.Println("platform", platform.Name)
	fmt.Println("platform", platform.Interfaces[1].Name)

	netController := utils.NewNetWorkController()
	tmp, _ := url.Parse("http://127.0.0.1:8080")
	netController.SetProxyURL(tmp)
	netController.SetInterface(platform.Interfaces[1])
	r := netController.Run()
	fmt.Println("response: ", r)

}

func TestReplaceMap(t *testing.T) {
	tests := models.ReplaceMap{Type: "headers", Key: "Cookies", Value: "test123", Interfaces: models.StringList{"TEST222", "TEST222"}}

	dbController := api.NewDBController()
	dbController.Connect("/Users/ryu/Documents/test.db")
	// var query map[string]interface{}
	platform := models.Platform{ID: 5, Name: "ceshi"}
	// dbController.CreatePlatform([]models.Platform{platform})
	platform, err := dbController.GetPlatform(platform)
	fmt.Println("err:", err)
	platform.ReplaceMaps = []models.ReplaceMap{tests}

	dbController.UpdatePlatform(platform)
	// dbController.DeletePlatforms([]models.Platform{platform})

}

func TestDB(t *testing.T) {
	dbController := api.NewDBController()
	dbController.Connect("/Users/ryu/Documents/test.db")

	interfaceInfo, err := dbController.GetInterface(models.Interface{Serial: "hEqbDdM2"})
	// fmt.Println(result)
	if err != nil {
		fmt.Println("接口获取错误：%w", err)
	}
	fmt.Println("interfaceInfo", interfaceInfo)
}

func TestRE(t *testing.T) {
	testStr := "设置一段文本![image1](images/test1)这是第二段文本![image2](images/test2)\n这是第三段文本![image2](images/test2)"
	reg := regexp.MustCompile(utils.ImgPattern)
	result := reg.FindAll([]byte(testStr), -1)
	s := reg.ReplaceAll([]byte(testStr), []byte("${2}AAAAA"))
	for _, item := range result {
		fmt.Println(string(item))
	}
	fmt.Println(string(s))
}

func appendQuery(queryList *[]string, queryParams *[]interface{}, querySQL string, value ...interface{}) {
	*queryList = append(*queryList, querySQL)
	*queryParams = append(*queryParams, value...)
}

func TestROd(t *testing.T) {

	// 设置浏览器
	browserPath, _ := launcher.LookPath()
	browserLauncher := launcher.New().Bin(browserPath).Headless(false)
	defer browserLauncher.Cleanup()
	browser := rod.New().NoDefaultDevice().ControlURL(browserLauncher.MustLaunch()).MustConnect().MustPage()

	// 访问登录页面
	browser.MustNavigate("https://passport.csdn.net/login?code=applets")

	for true {
		pages, _ := browser.Browser().Pages()
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println("len: ", len(pages))
		// 如果页面全部关闭，则推出
		if len(pages) == 0 {
			break

		} else {
			targetPage, err := pages.FindByURL("www.csdn.net")
			if err == nil {
				// 打印cookie
				cookies := targetPage.MustCookies()
				fmt.Println("cookies: ", cookies)
				for _, cookie := range cookies {
					fmt.Printf("Name: %s, Value: %s, Domain: %s\n", cookie.Name, cookie.Value, cookie.Domain)
				}
				break
			}

			time.Sleep(2 * time.Second)
		}
	}
	fmt.Println("done")
	// 关闭浏览器
	browserLauncher.Kill()

}
