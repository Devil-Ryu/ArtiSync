package test

import (
	"ArtiSync/backend/api"
	"ArtiSync/backend/models"
	"ArtiSync/backend/utils"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tidwall/gjson"
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
	controller.Connect("..//models/test.db")

	platform, _ := controller.GetPlatform(models.Platform{ID: 26})
	fmt.Println("platform", platform.Name)
	fmt.Println("platform", platform.Interfaces[2].Name)

	netController := utils.NewNetWorkController()
	netController.SetInterface(platform.Interfaces[2])
	r := netController.Run()
	fmt.Println("response: ", r)

}

func TestJSON(t *testing.T) {
	dir, err := os.UserConfigDir()
	if err != nil {
		println("获取用户配置目录失败: %w", err)
	}
	configPath := filepath.Join(dir, "ArtiSync")
	println("configPath", configPath)

	err = os.Mkdir(configPath, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("err1: ", err)
	}
}

func TestGJSON(t *testing.T) {

	test := `{"test":"ddddd","response": dasfdasf}`
	result := gjson.Parse(test)
	println(result.String())

	aaa := []string{}
	println(aaa == nil)
	fmt.Println(aaa)
}
