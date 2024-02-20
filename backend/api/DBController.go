package api

import (
	"ArtiSync/backend/models"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DBController 数据库控制器
type DBController struct {
	Ctx context.Context
	DB  *gorm.DB
}

// NewDBController 构造方法
func NewDBController() *DBController {
	return &DBController{}
}

// Startup 设置上下文并连接数据库
func (d *DBController) Startup(ctx context.Context) {
	d.Ctx = ctx
}

// Connect 数据库连接host
func (d *DBController) Connect(host string) (err error) {

	d.DB, err = gorm.Open(sqlite.Open(host), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("DBController Connect: %w", err)
	}
	err = d.AutoMigrate()
	if err != nil {
		return fmt.Errorf("DBController Migrate: %w", err)
	}
	return nil
}

// CheckConnect 检查数据库链接
func (d *DBController) CheckConnect() (err error) {
	if d.DB == nil {
		return errors.New("DBController not connect")
	}
	return nil
}

// AutoMigrate 迁移数据库
func (d *DBController) AutoMigrate() (err error) {
	// 检查数据库连接
	err = d.CheckConnect()
	if err != nil {
		return err
	}

	// 数据库迁移
	err = d.DB.AutoMigrate(&models.Platform{}, &models.Interface{}, &models.Header{}, &models.Body{}, &models.Query{}, &models.RPath{}, &models.ReplaceMap{}, &models.InterfaceRecord{})
	if err != nil {
		return err
	}
	return nil
}

// QueryInterfaceRecords 查询接口记录
func (d *DBController) QueryInterfaceRecords(query map[string]interface{}, pageNum int, pageSize int) (map[string]interface{}, error) {
	// 给页面大小赋初始值
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	queryList := []string{}
	queryParams := []interface{}{}

	for key, value := range query {
		// 判断该value值是否通过参数校验
		switch key {
		case "record_id":
			tmp, ok := value.(string)
			if ok && tmp != "" {
				appendQuery(&queryList, &queryParams, "record_id = ?", tmp)
			}
		case "date_time":
			tmp, ok := value.([]string)
			tmpList := []interface{}{}
			for _, item := range tmp {
				tmpList = append(tmpList, item)
			}
			if ok && len(tmp) != 0 {
				appendQuery(&queryList, &queryParams, "date_time between ? and ?", tmpList...)
			}
		case "platform_name":
			tmp, ok := value.([]string)
			if ok && len(tmp) != 0 {
				appendQuery(&queryList, &queryParams, "platform_name in ?", tmp)
			}
		case "status":
			tmp, ok := value.([]string)
			if ok && len(tmp) != 0 {
				appendQuery(&queryList, &queryParams, "status in ?", tmp)
			}
		case "tag":
			tmp, ok := value.(string)
			if ok && tmp != "" {
				appendQuery(&queryList, &queryParams, "tag = ?", tmp)
			}
		}

	}

	var result []models.InterfaceRecord

	querySQL := strings.Join(queryList, " and ")
	totalRows := d.DB.Where(querySQL, queryParams...).Find(&result).RowsAffected // 获取总行数
	d.DB.Where(querySQL, queryParams...).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&result)

	resultMap := map[string]interface{}{
		"result":    result,
		"pageNum":   pageNum,
		"pageSize":  pageSize,
		"totalRows": totalRows,
	}
	return resultMap, nil

}

// GetInterfaceRecords 获取接口记录
func (d *DBController) GetInterfaceRecords(query map[string]interface{}) (interfaceRecords []models.InterfaceRecord, err error) {
	// 检查数据库连接
	err = d.CheckConnect()
	if err != nil {
		return interfaceRecords, err
	}

	d.DB.Where(query).Preload(clause.Associations).Find(&interfaceRecords)

	return interfaceRecords, err

}

func appendQuery(queryList *[]string, queryParams *[]interface{}, querySQL string, value ...interface{}) {
	*queryList = append(*queryList, querySQL)
	*queryParams = append(*queryParams, value...)
}

// GetPlatforms 获取平台（批量）
func (d *DBController) GetPlatforms(query map[string]interface{}) (platforms []models.Platform, err error) {
	// 检查数据库连接
	err = d.CheckConnect()
	if err != nil {
		return platforms, err
	}

	d.DB.Where(query).Preload("Interfaces", func(db *gorm.DB) *gorm.DB { return db.Order("prior asc") }).
		Preload(clause.Associations).
		Preload("Interfaces.RequestURLPath").
		Preload("Interfaces.RequestQuery").
		Preload("Interfaces.RequestHeaders").
		Preload("Interfaces.RequestBody").
		Preload("Interfaces.Children", func(db *gorm.DB) *gorm.DB { return db.Order("prior asc") }).
		Preload("Interfaces.Children.RequestURLPath").
		Preload("Interfaces.Children.RequestQuery").
		Preload("Interfaces.Children.RequestHeaders").
		Preload("Interfaces.Children.RequestBody").
		Find(&platforms)
	return platforms, err
}

// GetPlatform 获取平台
func (d *DBController) GetPlatform(platform models.Platform) (models.Platform, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return platform, err
	}

	d.DB.Preload("Interfaces", func(db *gorm.DB) *gorm.DB { return db.Order("prior asc") }).
		Preload(clause.Associations).
		Preload("Interfaces.RequestURLPath").
		Preload("Interfaces.RequestQuery").
		Preload("Interfaces.RequestHeaders").
		Preload("Interfaces.RequestBody").
		Preload("Interfaces.Children", func(db *gorm.DB) *gorm.DB { return db.Order("prior asc") }).
		Preload("Interfaces.Children.RequestURLPath").
		Preload("Interfaces.Children.RequestQuery").
		Preload("Interfaces.Children.RequestHeaders").
		Preload("Interfaces.Children.RequestBody").
		Find(&platform)

	return platform, nil

}

// GetInterface 获取接口
func (d *DBController) GetInterface(interfaceInfo models.Interface) (models.Interface, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return interfaceInfo, err
	}
	d.DB.Preload(clause.Associations).
		Find(&interfaceInfo)
	return interfaceInfo, nil
}

// CreatePlatform 创建平台（可批量）
func (d *DBController) CreatePlatform(platforms []models.Platform) ([]models.Platform, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return platforms, err
	}
	d.DB.Create(platforms)
	return platforms, nil
}

// CreateInterface 创建接口
func (d *DBController) CreateInterface(interfaceInfo models.Interface) (models.Interface, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return interfaceInfo, err
	}
	d.DB.Create(&interfaceInfo)

	result, err := d.GetInterface(interfaceInfo)
	if err != nil {
		return interfaceInfo, err
	}
	return result, nil
}

// CreateOrUpdateInterfaceRecord 创建或更新接口记录
func (d *DBController) CreateOrUpdateInterfaceRecord(interfaceRecord models.InterfaceRecord) (models.InterfaceRecord, error) {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return interfaceRecord, err
	}
	d.DB.Save(&interfaceRecord)

	return interfaceRecord, nil
}

// UpdatePlatform 更新平台(不包括接口)
func (d *DBController) UpdatePlatform(platform models.Platform) error {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}
	// 添加新平台信息
	d.DB.Model(&models.Platform{ID: platform.ID}).Select("Name", "Disabled").Updates(platform)

	// 删除已有规则
	d.DB.Unscoped().Model(&models.Platform{ID: platform.ID}).Association("ReplaceMaps").Unscoped().Clear()
	// 更新新的规则
	d.DB.Unscoped().Model(&models.Platform{ID: platform.ID}).Association("ReplaceMaps").Replace(platform.ReplaceMaps)

	return nil

}

// UpdateInterface 更新接口
func (d *DBController) UpdateInterface(interfaces []models.Interface) error {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}
	// 删除原有接口及其对应参数
	d.DeleteInterfaces(interfaces)
	// 保存新的接口及其对应参数
	d.DB.Create(&interfaces)
	return nil

}

var count int = 0

// deleteInterfaceChildren 递归删除接口的children
func (d *DBController) deleteInterfaceChildren(interfaceInfo models.Interface) error {
	// 找到接口所有参数
	interfaceDetail, err := d.GetInterface(interfaceInfo)
	if err != nil {
		return err
	}
	// 如果有子接口
	if len(interfaceDetail.Children) != 0 {
		// 遍历子接口，并删除子接口
		for _, item := range interfaceDetail.Children {
			d.deleteInterfaceChildren(item)
		}
	}
	// 删除本接口
	d.DB.Select(clause.Associations).Delete(&interfaceDetail)
	return nil

}

// DeletePlatforms 批量删除平台
func (d *DBController) DeletePlatforms(platforms []models.Platform) error {
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}
	var interfaces []models.Interface
	// 找到平台的所有接口
	d.DB.Model(&platforms).Association("Interfaces").Find(&interfaces)

	// 删除接口及其其子接口对应参数
	for _, interfaceInfo := range interfaces {
		d.deleteInterfaceChildren(interfaceInfo)
	}

	// 删除平台及其对应关联
	d.DB.Select(clause.Associations).Delete(&platforms)
	return nil

}

// DeleteInterfaces 批量删除接口
func (d *DBController) DeleteInterfaces(interfaces []models.Interface) error {
	d.CheckConnect()
	// 检查数据库连接
	err := d.CheckConnect()
	if err != nil {
		return err
	}

	// 删除接口及其其子接口对应参数
	for _, interfaceInfo := range interfaces {
		fmt.Println("DeleteInterfaces", interfaceInfo)
		d.deleteInterfaceChildren(interfaceInfo)
	}
	return nil

}

// ExportPlatform 导出平台信息到文件
func (d *DBController) ExportPlatform(savePath string, platform models.Platform) (err error) {
	platformInfo, err := d.GetPlatform(platform)
	if err != nil {
		return err
	}
	// 先序列化
	platformJSON, err := json.Marshal(platformInfo)
	if err != nil {
		return fmt.Errorf("ExportPlatform  step1 marshal platform error: %w", err)
	}

	// 按照规定格式反序列化
	var platformFormat Platform
	err = json.Unmarshal(platformJSON, &platformFormat)
	if err != nil {
		return fmt.Errorf("ExportPlatform step2 unmarshal to format  error: %w", err)
	}

	// 拷贝platform为json的结构体
	platformFormatString, err := json.Marshal(platformFormat)
	if err != nil {
		return fmt.Errorf("ExportPlatform step3 json marshal error: %w", err)
	}

	// 格式化字符串
	var formatStr bytes.Buffer
	err = json.Indent(&formatStr, platformFormatString, "", "\t")
	if err != nil {
		return fmt.Errorf("ExportPlatform step4 format jsonStr error: %w", err)
	}

	// file, err := os.Create(platform.Name + ".json")
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("ExportPlatform create file error: %w", err)
	}
	bw := bufio.NewWriter(file)
	_, err = bw.WriteString(formatStr.String())
	if err != nil {
		return fmt.Errorf("ExportPlatform write content error: %w", err)
	}
	bw.Flush()
	return nil
}

// ImportPlatform 导入平台
func (d *DBController) ImportPlatform(fileDir string) (err error) {

	dataBytes, err := os.ReadFile(fileDir)
	if err != nil {
		return fmt.Errorf("ImportPlatform read file error: %w", err)
	}

	// 读取文件内容
	var platform models.Platform
	err = json.Unmarshal(dataBytes, &platform)
	if err != nil {
		return fmt.Errorf("ImportPlatform json marshal error: %w", err)
	}
	d.DB.Create(&platform)
	return nil
}

// GetConfigDir 获取用户配置目录
func (d *DBController) GetConfigDir() (dir string, err error) {
	systemType := runtime.GOOS
	if systemType == "darwin" {
		println("Platform: Mac")
	}

	// 获取用户配置目录
	dir, err = os.UserConfigDir()
	if err != nil {
		return dir, fmt.Errorf("获取用户配置目录失败: %w", err)
	}
	configPath := filepath.Join(dir, "ArtiSync")

	err = os.Mkdir(configPath, 0755)

	// 当不是文件存在错误的时候则抛出错误
	if err != nil && !os.IsExist(err) {
		fmt.Println("err1: ", err)
		return "", err
	}

	return configPath, nil
}

// GetConfigFilePath 获取配置文件
func (d *DBController) GetConfigFilePath() (configFilePath string, err error) {
	// 获取配置文件地址
	configPath, err := d.GetConfigDir()
	if err != nil {
		return configFilePath, err
	}
	configFilePath = filepath.Join(configPath, "config.json")

	defaultConfig := map[string]interface{}{
		"dbPath":      filepath.Join(configPath, "artisync.db"),
		"imagePath":   "",
		"imageSelect": "相对文章目录",
		"proxyURL":    "",
	}

	// 尝试加载配置文件
	_, err = d.LoadJSONFile(configFilePath)
	if err != nil {
		// 将默认配置进行保存
		err = d.SaveJSONFile(configFilePath, defaultConfig)
	}
	return configFilePath, err
}

// LoadJSONFile 加载JSON文件
func (d *DBController) LoadJSONFile(fileDir string) (jsonMap map[string]interface{}, err error) {
	dataBytes, err := os.ReadFile(fileDir)
	if err != nil {
		return jsonMap, fmt.Errorf("LoadJSONFile read file error: %w", err)
	}

	// 读取文件内容
	err = json.Unmarshal(dataBytes, &jsonMap)
	if err != nil {
		return jsonMap, fmt.Errorf("LoadJSONFile json Unmarshal error: %w", err)
	}
	return jsonMap, nil
}

// SaveJSONFile 加载JSON文件
func (d *DBController) SaveJSONFile(fileDir string, jsonMap map[string]interface{}) (err error) {
	// 先序列化
	jsonStr, err := json.Marshal(jsonMap)
	if err != nil {
		return fmt.Errorf("SaveJSONFile  marshal platform error: %w", err)
	}

	// 格式化字符串
	var formatStr bytes.Buffer
	err = json.Indent(&formatStr, jsonStr, "", "\t")
	if err != nil {
		return fmt.Errorf("SaveJSONFile format jsonStr error: %w", err)
	}

	file, err := os.Create(fileDir)
	if err != nil {
		return fmt.Errorf("SaveJSONFile create file error: %w", err)
	}

	bw := bufio.NewWriter(file)
	_, err = bw.WriteString(formatStr.String())
	if err != nil {
		return fmt.Errorf("SaveJSONFile write content error: %w", err)
	}
	bw.Flush()
	return nil
}
