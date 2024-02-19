package utils

// ResponseJSON 响应
type ResponseJSON struct {
	StatusCode int         // 状态码
	Message    string      // 返回信息
	Data       interface{} // 返回数据
}

// 状态枚举
const (
	Publishing       string = "发布中"  // 发布中（文章）
	PublishWating    string = "待发布"  // 待发布（文章）
	PublishedSuccess string = "发布成功" // 发布成功 （文章）
	PublishedFailed  string = "发布失败" // 发布失败（文章）
	RunningSuccess   string = "运行成功" // 运行成功（接口）
	RunningFailed    string = "运行失败" // 运行失败（接口）
	Waiting          string = "等待中"  // 等待中
	Running          string = "运行中"  // 运行中
)

// // CommonUtils 常用工具
// type CommonUtils struct {
// 	Ctx context.Context
// }

// // GetConfigDir 获取用户配置目录
// func (c *CommonUtils) GetConfigDir() (dir string, err error) {
// 	systemType := runtime.GOOS
// 	if systemType == "darwin" {
// 		println("Platform: Mac")
// 	}

// 	// 获取用户配置目录
// 	dir, err = os.UserConfigDir()
// 	if err != nil {
// 		return dir, fmt.Errorf("获取用户配置目录失败: %w", err)
// 	}
// 	configPath := filepath.Join(dir, "ArtiSync")

// 	err = os.Mkdir(configPath, 0755)

// 	// 当不是文件存在错误的时候则抛出错误
// 	if err != nil && !os.IsExist(err) {
// 		fmt.Println("err1: ", err)
// 		return "", err
// 	}

// 	return configPath, err
// }

// // LoadJSONFile 加载JSON文件
// func (c *CommonUtils) LoadJSONFile(fileDir string) (jsonMap map[string]interface{}, err error) {
// 	dataBytes, err := os.ReadFile(fileDir)
// 	if err != nil {
// 		return jsonMap, fmt.Errorf("LoadJSONFile read file error: %w", err)
// 	}

// 	// 读取文件内容
// 	err = json.Unmarshal(dataBytes, &jsonMap)
// 	if err != nil {
// 		return jsonMap, fmt.Errorf("LoadJSONFile json Unmarshal error: %w", err)
// 	}
// 	return jsonMap, nil
// }

// // SaveJSONFile 加载JSON文件
// func (c *CommonUtils) SaveJSONFile(fileDir string, jsonMap map[string]interface{}) (err error) {
// 	// 先序列化

// 	jsonStr, err := json.Marshal(jsonMap)
// 	if err != nil {
// 		return fmt.Errorf("SaveJSONFile  marshal platform error: %w", err)
// 	}

// 	// 格式化字符串
// 	var formatStr bytes.Buffer
// 	err = json.Indent(&formatStr, jsonStr, "", "\t")
// 	if err != nil {
// 		return fmt.Errorf("SaveJSONFile format jsonStr error: %w", err)
// 	}

// 	file, err := os.Create(fileDir)
// 	if err != nil {
// 		return fmt.Errorf("SaveJSONFile create file error: %w", err)
// 	}
// 	bw := bufio.NewWriter(file)
// 	_, err = bw.WriteString(formatStr.String())
// 	if err != nil {
// 		return fmt.Errorf("SaveJSONFile write content error: %w", err)
// 	}
// 	bw.Flush()
// 	return nil
// }
