package test

import (
	"ArtiSync/backend/models"
	"ArtiSync/backend/utils"
	"fmt"
	"regexp"
	"testing"
)

func TestCase(t *testing.T) {

	testInterface := models.Interface{
		Name:            "测试接口",
		Serial:          "TEST1111",
		RequestMethod:   "POST",
		RequestURL:      "https://test.example.com/",
		RequestHeaders:  []models.Header{{Key: "test1", Value: "test11"}, {Key: "test2", Value: "test22"}},
		RequestQuery:    []models.Query{{Key: "test1", Value: "test11"}, {Key: "test2", Value: "test22"}},
		RequestBodyType: "JSON",
		RequestBody:     []models.Body{{Key: "body1", Value: "body11"}, {Key: "body2", Value: "body22"}},
	}

	netController := utils.NewNetWorkController()

	netController.SetInterface(testInterface)
	netController.Run()
	response, err := netController.GetResponseMappedValue()
	if err != nil {
		fmt.Println("错误信息", err.Error())
	}
	fmt.Println("response", response)
	// fmt.Println(netController.ResponsePool)

}

func TestCustom(t *testing.T) {
	re, err := regexp.Compile(`\d+`)
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		a := re.FindString("test")
		fmt.Println("a: ", a)
	}

}
