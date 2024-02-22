package test

import (
	"ArtiSync/backend/api"
	"ArtiSync/backend/models"
	"ArtiSync/backend/utils"
	"fmt"
	"net/url"
	"strings"
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

	// tmp := map[string]interface{}{
	// 	"record_id":     "VGIs83uE",
	// 	"platform_name": "MYBLOG",
	// 	"serial":        "sW5OnfAB",
	// }
	// test, _ := dbController.GetInterfaceRecords(tmp)

	// fmt.Println(test)
	query := map[string]interface{}{
		"record_id":     "8lX4oeOV",
		"date_time":     []string{"2024-02-18 00:00:00", "2024-02-20 23:59:59"},
		"platform_name": []string{"MYBLOG"},
	}

	queryList := []string{}
	queryParams := []interface{}{}

	for key, value := range query {
		// 判断该value值是否通过参数校验
		switch key {
		case "record_id":
			tmp, ok := value.(string)
			if ok {
				appendQuery(&queryList, &queryParams, "record_id = ?", tmp)
			}
		case "date_time":
			tmp, ok := value.([]string)
			tmpList := []interface{}{}
			for _, item := range tmp {
				tmpList = append(tmpList, item)
			}
			if ok {
				appendQuery(&queryList, &queryParams, "date_time between ? and ?", tmpList...)
			}
		case "platform_name":
			tmp, ok := value.([]string)
			if ok {
				appendQuery(&queryList, &queryParams, "platform_name in ?", tmp)
			}
		}
	}

	querySQL := strings.Join(queryList, " and ")
	fmt.Println(querySQL)
	fmt.Println(queryParams)

	var result []models.InterfaceRecord

	dbController.DB.Where(querySQL, queryParams...).Find(&result)
	fmt.Println(len(result))
	// fmt.Println(result)
}

func appendQuery(queryList *[]string, queryParams *[]interface{}, querySQL string, value ...interface{}) {
	*queryList = append(*queryList, querySQL)
	*queryParams = append(*queryParams, value...)
}
