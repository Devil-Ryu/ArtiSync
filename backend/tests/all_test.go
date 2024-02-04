package test

import (
	"ArtiSync/backend/api"
	"ArtiSync/backend/models"
	"ArtiSync/backend/utils"
	"fmt"
	"testing"
	"time"
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
