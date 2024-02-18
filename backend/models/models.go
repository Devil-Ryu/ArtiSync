package models

import (
	"database/sql/driver"
	"math/rand"
	"strings"

	"gorm.io/gorm"
)

// StringList 字符列表
type StringList []string

// Platform 平台配置
type Platform struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       // 平台名称
	Disabled    bool         `gorm:"default:false"` // 是否禁用
	Interfaces  []Interface  // 平台下的接口（一对多关系）
	ReplaceMaps []ReplaceMap // 平台下的替换规则
}

// Interface 接口（组）配置
type Interface struct {
	ID               uint        `gorm:"primaryKey"`
	PlatformID       uint        // 平台ID
	Prior            uint        // 接口优先级
	Serial           string      // 接口编号, 一个平台内接口的内部编号，平台内编号不能重
	Name             string      // 接口名称
	IsGroup          bool        `gorm:"default:false"` // 是否为组
	ParentSerial     *string     // 父接口ID
	Children         []Interface `gorm:"foreignkey:ParentSerial;references:Serial"` // 子任务
	Type             string      `gorm:"default:normal"`                            // 接口类型(normal，images， article，front)
	RequestURL       string      //  接口地址（入参）
	RequestURLPath   []RPath     // 接口地址拼接
	RequestMethod    string      // 请求方法（入参）
	RequestQuery     []Query     // 请求查询参数（入参）
	RequestHeaders   []Header    // 请求头（入参）
	RequestBody      []Body      // 请求体（入参）
	RequestBodyType  string      `gorm:"default:NONE"` // 请求体格式（JSON-字符串、ROWDATA-字符串、FORMDATA-表单、NONE-无）
	ResponseType     string      `gorm:"default:NONE"` // 响应体格式（JSON、TEXT、NONE）（出参）
	ResponseTemple   string      // 响应体模板（JSON、RE、NONE）用于解析响应体（出参）
	ResponseValidate bool        `gorm:"default:false"` // 是否校验校验响应包，若为true则接口发送后会按照响应体规则进行校验
	ResponseMapType  string      `gorm:"default:NONE"`  // 响应体映射是否动态，该接口响应根据映射出参（JSON-返回匹配包、RE-返回匹配包、NONE-返回原始包）用于解析响应体（出参）
	ResponseMapRule  string      // 响应体映射，该接口响应根据映射出参（JSON、RE、NONE）用于解析响应体（出参）

}

// RPath 请求路径
type RPath struct {
	ID          uint `gorm:"primaryKey"`
	InterfaceID uint
	Dynamic     bool   `gorm:"default:false"` // 参数是否动态获取（默认为false,如果为其他的，会修改对应位置的值）
	Value       string // 值
}

// Header 请求头
type Header struct {
	ID          uint `gorm:"primaryKey"`
	InterfaceID uint
	Convert     bool   `gorm:"default:false"`
	Dynamic     bool   `gorm:"default:false"` // 参数是否动态获取（默认为false,如果为其他的，会修改对应位置的值）
	Key         string // 键
	Value       string // 值
}

// ReplaceMap 替换Map
type ReplaceMap struct {
	ID         uint `gorm:"primaryKey"`
	PlatformID uint
	Disabled   bool       `gorm:"default:false"`  // 是否禁用
	Type       string     `gorm:"default:header"` // 替换的位置
	Key        string     // 键
	Value      string     // 值
	Interfaces StringList // 应用的接口编号
}

// Query 请求参数
type Query struct {
	ID          uint `gorm:"primaryKey"`
	InterfaceID uint
	Dynamic     bool   `gorm:"default:false"` // 参数是否动态获取（默认为false,如果为其他的，会修改对应位置的值）
	Key         string // 键
	Value       string // 值
}

// Body 请求体
type Body struct {
	ID          uint `gorm:"primaryKey"`
	InterfaceID uint
	IsFile      bool   `gorm:"default:false"` // 是否为文件
	FileName    string // 文件名称
	Convert     bool   `gorm:"default:false"`
	Dynamic     bool   `gorm:"default:false"` // 参数是否动态获取（默认为false,如果为其他的，会修改对应位置的值）
	Key         string // 键
	Value       string // 值
}

// InterfaceRecord 接口运行记录
type InterfaceRecord struct {
	ID              uint `gorm:"primaryKey"`
	DateTime        string
	RecordID        string // 记录ID（每次运行生成一个记录ID，用于分辨是那一次批量运行的记录）
	ArticleName     string // 文章名称
	PlatformName    string // 平台名称
	Serial          string // 接口编号
	Name            string // 接口名称
	RequestURL      string //  接口地址（入参）
	RequestMessage  string `gorm:"type:text"` // 请求报文
	ResponseMessage string `gorm:"type:text"` // 响应报文
	Status          string // 运行状态
	Tag             string // 标签

}

// BeforeCreate 接口创建前
func (i *Interface) BeforeCreate(tx *gorm.DB) (err error) {
	if i.Serial == "" {
		i.Serial = GenerateRandomKey(8)
	}
	return
}

// Scan 查询方法
func (s *StringList) Scan(valueStr interface{}) error {
	tmp := valueStr.(string)
	valueArr := strings.Split(tmp, "|")
	*s = valueArr
	return nil
}

// Value 储存方法
func (s StringList) Value() (driver.Value, error) {
	str := strings.Join(s, "|")
	return str, nil
}

// // BeforeDelete Hook 资源级联删除接口记录
//
//	func (resource *Resource) BeforeDelete(tx *gorm.DB) (err error) {
//		// 查询资源对应的接口
//		var interfaces []Interface
//		tx.Model(resource).Association("InterfacesRequest").Find(&interfaces)
//
//		// 遍历所有接口，接触资源与接口之间的关联
//		for _, interface_ := range interfaces {
//			tx.Model(resource).Association("InterfacesRequest").Delete(interface_)
//		}
//
//		return
//	}
//
// // BeforeDelete Hook 参数级联删除接口记录
//
//	func (parameter *Parameter) BeforeDelete(tx *gorm.DB) (err error) {
//		// 查询参数对应的接口
//		var interfaces []Interface
//
//		// 查看该参数是否和InterfacesRequest有关联
//		tx.Model(parameter).Association("InterfacesRequest").Find(&interfaces)
//		// 遍历所有接口，接触参数与接口之间的关联
//		for _, interface_ := range interfaces {
//			tx.Model(parameter).Association("InterfacesRequest").Delete(interface_)
//		}
//
//		// 查看该参数是否和InterfacesResponse有关联
//		tx.Model(parameter).Association("InterfacesResponse").Find(&interfaces)
//		// 遍历所有接口，接触参数与接口之间的关联
//		for _, interface_ := range interfaces {
//			tx.Model(parameter).Association("InterfacesResponse").Delete(interface_)
//		}
//		return
//	}
func GenerateRandomKey(length int) string {
	// 可以根据需要修改字符集
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var result string
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSet))
		result += string(charSet[randomIndex])
	}
	return result
}

/*
ReplaceMaps: [{"GroupName": "xxxx1", "Maps": []}, ]
*/
