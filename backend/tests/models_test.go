package test

import (
	"ArtiSync/backend/models"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestCreateModels(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("/backend/models/test.db"), &gorm.Config{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&models.Platform{}, &models.Interface{}, &models.Header{}, &models.Body{}, &models.Query{}, &models.RPath{})
	if err != nil {
		return
	}
}

func TestCreateData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return
	}

	platform1 := &models.Platform{Name: "CSDN",
		Interfaces: []models.Interface{
			{
				Name:          "接口1",
				RequestURL:    "https://www.interface1.com.cn",
				RequestMethod: "GET",
				RequestQuery: []models.Query{
					{Key: "query1-1", Value: "query1-11"},
					{Key: "query1-2", Value: "query1-22"},
					{Key: "query1-3", Value: `1|resp2|resp1`, Dynamic: true},
				},
				RequestHeaders: []models.Header{
					{Key: "header1-1", Value: "header1-11"},
					{Key: "header1-2", Value: "header1-22"},
					{Key: "header1-3", Value: `1|resp1|resp2`, Dynamic: true},
				},
				RequestBody: []models.Body{
					{Key: "body1-1", Value: "body1-11"},
					{Key: "body1-2", Value: "body1-22"},
					{Key: "body1-3", Value: `1|data|src`, Dynamic: true},
				},
				ResponseType:   "JSON",
				ResponseTemple: `{"resp1":{"resp1":"resp1Value"}}`,
			},
			{
				Name:          "接口2",
				RequestURL:    "https://www.interface2.com.cn",
				RequestMethod: "POST",
				RequestQuery: []models.Query{
					{Key: "query2-1", Value: "query2-11"},
					{Key: "query2-2", Value: "query2-22"},
					{Key: "query2-3", Value: `2|resp2|resp3`, Dynamic: true},
				},
				RequestHeaders: []models.Header{
					{Key: "header2-1", Value: "header2-11"},
					{Key: "header2-2", Value: "header2-22"},
					{Key: "header2-3", Value: `2|resp2|resp2`, Dynamic: true},
				},
				RequestBody: []models.Body{
					{Key: "body2-1", Value: "body2-11"},
					{Key: "body2-2", Value: "body2-22"},
					{Key: "body2-3", Value: `2|resp2|resp2`, Dynamic: true},
				},
				ResponseType:   "JSON",
				ResponseTemple: `{"resp2":{"resp2":"resp2Value"}}`,
			},
		},
	}

	platform2 := &models.Platform{Name: "MYBLOG",
		Interfaces: []models.Interface{
			{
				Name:          "接口3",
				RequestURL:    "https://www.interface3.com.cn",
				RequestMethod: "PUT",
				RequestQuery: []models.Query{
					{Key: "query3-1", Value: "query3-11"},
					{Key: "query3-2", Value: "query3-22"},
					{Key: "query3-3", Value: `1|resp3|resp3`, Dynamic: true},
				},
				RequestHeaders: []models.Header{
					{Key: "header3-1", Value: "header3-11"},
					{Key: "header3-2", Value: "header3-22"},
					{Key: "header3-3", Value: `1|resp3|resp3`, Dynamic: true},
				},
				RequestBody: []models.Body{
					{Key: "body3-1", Value: "body3-11"},
					{Key: "body3-2", Value: "body3-22"},
					{Key: "body3-3", Value: `1|resp3|resp3`, Dynamic: true},
				},
				ResponseType:   "JSON",
				ResponseTemple: `{"resp3":{"resp3":"resp3Value"}}`,
			},
		},
	}

	//// 批量创建平台数据
	db.Create(platform1)
	db.Create(platform2)

}

func TestResetDB(t *testing.T) {

	// 重新创建模型
	TestCreateModels(t)
	// 重新更新数据
	TestCreateData(t)

}

func TestSQL(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return
	}

	platform := map[string]interface{}{
		"Name": "平台1",
		"Interfaces": `
			{
				"RequestURL":    "https://wwww.test.com",
				"RequestMethod": "POST"
			}`,
	}

	db.Model(&models.Platform{}).Create(platform)

}

func TestGetInterface(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("../models/test.db"), &gorm.Config{})
	if err != nil {
		return
	}

	var interfaces []models.Interface
	db.Model(&models.Interface{}).Preload(clause.Associations).Find(&interfaces)
	fmt.Println(interfaces)
}
