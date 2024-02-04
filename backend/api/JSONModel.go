package api

// 镜像复制models/models.go中的文件，用于自定义JSON序列化内容时用

// Platform 平台配置
type Platform struct {
	ID          uint         `json:"-"`
	Name        string       // 平台名称
	Disabled    bool         // 是否禁用
	Interfaces  []Interface  // 平台下的接口（一对多关系）
	ReplaceMaps []ReplaceMap // 平台下的替换规则
}

// Interface 接口（组）配置
type Interface struct {
	ID               uint        `json:"-"`
	PlatformID       uint        `json:"-"` // 平台ID
	Prior            uint        // 接口优先级
	Serial           string      // 接口编号, 一个平台内接口的内部编号，平台内编号不能重
	Name             string      // 接口名称
	IsGroup          bool        // 是否为组
	ParentSerial     *string     // 父接口ID
	Children         []Interface // 子任务
	Type             string      // 接口类型(normal，images， article，front)
	RequestURL       string      //  接口地址（入参）
	RequestURLPath   []RPath     // 接口地址拼接
	RequestMethod    string      // 请求方法（入参）
	RequestQuery     []Query     // 请求查询参数（入参）
	RequestHeaders   []Header    // 请求头（入参）
	RequestBody      []Body      // 请求体（入参）
	RequestBodyType  string      // 请求体格式（JSON-字符串、ROWDATA-字符串、FORMDATA-表单、NONE-无））
	ResponseType     string      // 响应体格式（JSON、TEXT、NONE）（出参）
	ResponseTemple   string      // 响应体内容（JSON、re、NONE）用于解析响应体（出参）
	ResponseValidate bool        // 是否校验校验响应包，若为true则接口发送后会按照响应体规则进行校验
	ResponseMapType  string      // 响应体映射是否动态，该接口响应根据映射出参（JSON、re、NONE）用于解析响应体（出参）
	ResponseMapRule  string      // 响应体映射，该接口响应根据映射出参（JSON、re、NONE）用于解析响应体（出参）

}

// RPath 请求路径
type RPath struct {
	ID          uint   `json:"-"`
	InterfaceID uint   `json:"-"`
	Dynamic     bool   // 参数是否动态获取（默认为false,如果为其他的，会修改对应位置的值）
	Value       string // 值
}

// ReplaceMap 替换Map
type ReplaceMap struct {
	ID         uint     `json:"-"`
	PlatformID uint     `json:"-"`
	Disabled   bool     `gorm:"default:false"`  // 是否禁用
	Type       string   `gorm:"default:header"` // 替换的位置
	Key        string   // 键
	Value      string   // 值
	Interfaces []string // 应用的接口编号
}

// Header 请求头
type Header struct {
	ID          uint   `json:"-"`
	InterfaceID uint   `json:"-"`
	Dynamic     bool   // 参数是否动态获取（默认为false,如果为其他的，会修改对应位置的值）
	Key         string // 键
	Value       string // 值
}

// Query 请求参数
type Query struct {
	ID          uint   `json:"-"`
	InterfaceID uint   `json:"-"`
	Dynamic     bool   // 参数是否动态获取（默认为false,如果为其他的，会修改对应位置的值）
	Key         string // 键
	Value       string // 值
}

// Body 请求体
type Body struct {
	ID          uint   `json:"-"`
	InterfaceID uint   `json:"-"`
	IsFile      bool   // 是否为文件
	FileName    string // 文件名称
	Dynamic     bool   // 参数是否动态获取（默认为false,如果为其他的，会修改对应位置的值）
	Key         string // 键
	Value       string // 值
}
