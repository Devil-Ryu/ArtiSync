package utils

import (
	"ArtiSync/backend/models"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// NetWorkController 网络控制器(TODO 增加睡眠时长，以及发送失败重试次数)
type NetWorkController struct {
	Ctx                context.Context
	ProxURL            *url.URL
	CurRequestMessage  string
	CurResponseMessage string
	// URL          *url.URL
	// Request      *http.Request
	Interface        models.Interface
	ResponsePool     map[string]string
	ResponsePoolType map[string]string
}

// NewNetWorkController 生成网络控制器
func NewNetWorkController() *NetWorkController {
	return &NetWorkController{
		ResponsePool:     make(map[string]string),
		ResponsePoolType: make(map[string]string),
	}
}

// StartUp  wails调用方法，用于传入上下文
func (n *NetWorkController) StartUp(ctx context.Context) {
	n.Ctx = ctx
}

// ResetResponsePool 重置响应池
func (n *NetWorkController) ResetResponsePool() {
	n.ResponsePool = make(map[string]string)
	n.ResponsePoolType = make(map[string]string)
}

// SetInterface 设置接口
func (n *NetWorkController) SetInterface(interfaceInfo models.Interface) {
	n.Interface = interfaceInfo
}

// SetProxyURL 设置代理
func (n *NetWorkController) SetProxyURL(proxURL *url.URL) {
	n.ProxURL = proxURL
}

// base64encode base64编码
func (n *NetWorkController) base64encode(srcBytes []byte) string {
	encodeStr := base64.StdEncoding.EncodeToString(srcBytes)
	return encodeStr
}

// base64encode base64解码
func (n *NetWorkController) base64decode(encodeStr string) string {
	decodeStr, err := base64.StdEncoding.DecodeString(encodeStr)
	if err != nil {
		fmt.Println("err: ", err.Error())
	}

	return string(decodeStr)
}

func (n *NetWorkController) convertValue(inputValue string) (output interface{}) {
	// 尝试转化为整型
	output, err := strconv.Atoi(inputValue)
	if err == nil {
		return output
	}

	// 尝试转化为浮点型
	output, err = strconv.ParseFloat(inputValue, 64)
	if err == nil {
		return output
	}
	return inputValue
}

// GetDynamicParam 获取动态参数(修改为判断对应接口的返回值类型)
func (n *NetWorkController) GetDynamicParam(paramValue string) (retVal string, err error) {
	paramArr := strings.Split(paramValue, ".") // 转换动态参数为数组
	var responseType string
	responseType = n.ResponsePoolType[paramArr[0]] // 找到对应接口的类型
	if responseType == "" {
		return retVal, fmt.Errorf("网络请求池未找到响应类型[%s]", paramArr[0])
	}

	retVal = n.ResponsePool[paramArr[0]] // 默认直接返回响应池中的值
	// 如果是json，根据value[0]找到响应包，并把response转化为json，然后用gjson去到对应的值，根据responseType 解析paramValue 1.data.src
	if responseType == "JSON" {
		responseJSON := gjson.Parse(n.ResponsePool[paramArr[0]]) // 在响应池中根据接口ID找到对应接口响应，转化为GJSON格式
		if !responseJSON.Exists() {
			return retVal, fmt.Errorf("解析响应包至JSON失败")
		}
		if len(paramArr) > 1 {
			paramPath := strings.Join(paramArr[1:], ".")
			result := responseJSON.Get(paramPath) // 根据参数路径找到对应响应的值
			if !result.Exists() {
				fmt.Println("err", err)
				return retVal, fmt.Errorf("匹配JSON响应失败: %w, 响应包: %s, 匹配路径: %s", err, n.ResponsePool[paramArr[0]], paramPath)
			}
			retVal = result.String()
		} else {
			retVal = responseJSON.String() // 如果参数子项小于1则返回整体报文
			if err != nil {
				fmt.Println("err", err)
				return retVal, fmt.Errorf("匹配JSON响应失败(路径参数为长度为0): %w", err)
			}
		}

		fmt.Println("[*] 解析动态参数：", paramValue)
		// fmt.Println("[*] 找到的响应：", n.ResponsePool[paramArr[0]])
		// fmt.Println("[*] 参数路径：", paramPath)
		fmt.Println("[*] 参数长度：", len(paramArr))
		fmt.Println("[*] 返回值长度： ", len(retVal))

		if !strings.HasPrefix(retVal, `"`) && !strings.HasPrefix(retVal, `'`) {
			retVal = strconv.Quote(retVal) // 加上引号，处理部分返回不是字符串的情况
		}
		retVal, _ = strconv.Unquote(retVal) // 去除转义字符和引号，因为json序列化的时候会转义
	}

	// if len(retVal) > 20 {
	// 	fmt.Println("[*] 返回值>20： ", retVal[:10])
	// } else {
	// 	fmt.Println("[*] 返回值<20： ", retVal)
	// }

	return retVal, nil
}

// GetResponseMappedValue 获取响应体（TODO 添加ResponseMapType判断）
func (n *NetWorkController) GetResponseMappedValue() (retVal string, err error) {

	// 响应体格式为NONE，则返回完整报文
	if n.Interface.ResponseMapType == "NONE" {
		retVal = n.ResponsePool[fmt.Sprint(n.Interface.Serial)]
		return retVal, nil
	}

	// 响应体格式RE，则进行正则匹配，获取不到则返回空
	if n.Interface.ResponseMapType == "RE" {
		re, err := regexp.Compile(n.Interface.ResponseMapRule)
		if err != nil {
			return retVal, err
		}
		retVal = n.ResponsePool[fmt.Sprint(n.Interface.Serial)]
		retVal = re.FindString(retVal) // 根据正则表达式匹配

		return retVal, nil
	}

	// 响应体格式JSON，则动态取值，JSON默认动态取值，获取不到则返回空
	if n.Interface.ResponseMapType == "JSON" {
		retVal, err = n.GetDynamicParam(n.Interface.ResponseMapRule)
		if err != nil {
			return retVal, err
		}
		return retVal, nil
	}

	// 不属于上述情况，则直接获取
	retVal = n.ResponsePool[fmt.Sprint(n.Interface.Serial)]
	return retVal, nil
}

// ParseBody 解析请求体
func (n *NetWorkController) ParseBody(bodyType string, bodyData []models.Body) (requestData string, contentType string, err error) {

	// 请求体为JSON时（只支持一级，不支持多级JSON）
	if strings.ToUpper(bodyType) == "JSON" {
		jsonData := map[string]interface{}{}
		for _, item := range bodyData {
			tmp := item.Value
			if item.Dynamic {
				itemValue, err := n.GetDynamicParam(item.Value)
				if err != nil {
					return requestData, contentType, fmt.Errorf("Err|ParseBody: %w", err)
				}
				jsonData[item.Key] = itemValue
				tmp = itemValue
			} else {
				jsonData[item.Key] = item.Value
			}
			if item.Convert {
				jsonData[item.Key] = n.convertValue(tmp)
			}
		}
		jsonStr, _ := json.Marshal(jsonData)
		requestData = string(jsonStr)
		contentType = "application/json"
		// fmt.Println("requestBody(JSON)")
		// fmt.Println(requestData)
	}

	// 请求体为JSONTEXT时（纯Text则只取第一个）
	if strings.ToUpper(bodyType) == "JSONTEXT" {
		if bodyData[0].Dynamic {
			itemValue, err := n.GetDynamicParam(bodyData[0].Value)
			if err != nil {
				return requestData, contentType, fmt.Errorf("Err|ParseBody: %w", err)
			}
			requestData = itemValue
		} else {
			requestData = bodyData[0].Value
		}
		contentType = "application/json"
		// fmt.Println("requestBody(JSON)")
		// fmt.Println(requestData)

	}

	// rowData数据格式：eg &test=1&aaa=2
	if strings.ToUpper(bodyType) == "ROWDATA" {
		formValues := url.Values{}
		for _, item := range bodyData {
			if item.Dynamic {
				itemValue, err := n.GetDynamicParam(item.Value)
				if err != nil {
					return requestData, contentType, fmt.Errorf("Err|ParseBody: %w", err)
				}
				formValues.Set(item.Key, itemValue)
			} else {
				formValues.Set(item.Key, item.Value)
			}
		}
		requestData = formValues.Encode()
		contentType = "application/x-www-form-urlencoded"
		// fmt.Println("requestBody(ROWDATA)")
		// fmt.Println(requestData)
	}

	// FORMDATA 表单（一般传文件用这个格式）
	if strings.ToUpper(bodyType) == "FORMDATA" {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for _, item := range bodyData {
			if item.Dynamic {
				// 动态参数逻辑
				itemValue, err := n.GetDynamicParam(item.Value) // 获取动态参数
				if err != nil {
					return requestData, contentType, fmt.Errorf("解析BODY错误: %w", err)
				}

				// 判断文件名是否为动态
				fileName := item.FileName // 默认文件名为填入的名称
				if item.FileNameDynamic {
					fileName, err = n.GetDynamicParam(item.FileName) // 获取动态参数
					if err != nil {
						return requestData, contentType, fmt.Errorf("解析BODY错误: %w", err)
					}
				}

				if item.IsFile {
					part, _ := writer.CreateFormFile(item.Key, fileName) // 创建文件表单
					part.Write([]byte(itemValue))                        // 写入文件内容
					body.Write([]byte(fmt.Sprintf("\r\n--%s--\r\n", writer.Boundary())))
				} else {
					writer.WriteField(item.Key, itemValue) // 写入非文件内容
				}
			} else {
				// 非动态参数逻辑

				// 判断文件名是否为动态
				fileName := item.FileName // 默认文件名为填入的名称
				if item.FileNameDynamic {
					fileName, err = n.GetDynamicParam(item.FileName) // 获取动态参数
					if err != nil {
						return requestData, contentType, fmt.Errorf("解析BODY错误: %w", err)
					}
				}

				if item.IsFile {
					part, _ := writer.CreateFormFile(item.Key, fileName) // 创建文件表单
					part.Write([]byte(item.Value))                       // 写入文件内容
					body.Write([]byte(fmt.Sprintf("\r\n--%s--\r\n", writer.Boundary())))
				} else {
					writer.WriteField(item.Key, item.Value) // 写入非文件内容
				}
			}
		}
		// 添加分隔符结束标志
		// fmt.Println("requestBody(FORMDATA)")
		// fmt.Println(body)
		requestData = body.String()
		contentType = writer.FormDataContentType()
	}

	return requestData, contentType, nil
}

// GetRequest 整体流程，1.读取参数 2.参数设置 3.发送请求 4.返回结果
func (n *NetWorkController) GetRequest() (request *http.Request, err error) {

	// 1.设置URL
	RequestURL, err := url.ParseRequestURI(n.Interface.RequestURL)
	if err != nil {
		return request, fmt.Errorf("解析URL错误: %w", err)
	}

	// 2.添加URL路径
	for _, rPath := range n.Interface.RequestURLPath {
		if rPath.Dynamic {
			itemValue, err := n.GetDynamicParam(rPath.Value)
			if err != nil {
				return request, fmt.Errorf("获取动态URL路径错误[%s]: %w", rPath.Value, err)
			}
			RequestURL.Path = path.Join(RequestURL.Path, "/"+itemValue)
		} else {
			RequestURL.Path = path.Join(RequestURL.Path, "/"+rPath.Value)
		}
	}

	// 3.设置query参数
	urlData := url.Values{}
	for _, item := range n.Interface.RequestQuery {
		if item.Dynamic {
			itemValue, err := n.GetDynamicParam(item.Value)
			if err != nil {
				return request, fmt.Errorf("获取动态查询参数错误[%s]: %w", item.Key, err)
			}
			urlData.Add(item.Key, itemValue)
		} else {
			urlData.Add(item.Key, item.Value)
		}
	}

	RequestURL.RawQuery = urlData.Encode() // 参数更新到URL中

	// 4.设置body参数()
	bodyData, contentType, err := n.ParseBody(n.Interface.RequestBodyType, n.Interface.RequestBody)
	if err != nil {
		return request, fmt.Errorf("解析请求体错误: %w", err)
	}

	// 5.创建Request对象

	request, _ = http.NewRequest(n.Interface.RequestMethod, RequestURL.String(), bytes.NewReader([]byte(bodyData)))

	// 6.设置Headers参数
	// 设置contentType
	for _, header := range n.Interface.RequestHeaders {
		if header.Dynamic {
			headerValue, err := n.GetDynamicParam(header.Value)
			if err != nil {
				return request, fmt.Errorf("获取动态Headers参数错误[%s]: %w", header.Key, err)
			}
			request.Header.Add(header.Key, headerValue)
		} else {
			request.Header.Add(header.Key, header.Value)
		}
	}
	request.Header.Set("Content-Type", contentType)

	return request, nil

}

// storeToPool 存进储存池
func (n *NetWorkController) storeToPool(key string, value []byte, valueType string) {
	valueType = strings.ToUpper(valueType) // 统一转化为大写

	// // 处理存储值为JSON的情况
	// if valueType == "JSON" {
	// 	valueStr := string(value)  // 先将value转化为字符串
	// 	gjson

	// }
}

// Run 网络请求模块运行方法
func (n *NetWorkController) Run() ResponseJSON {
	// time.Sleep(1 * time.Second) // 睡眠1s
	transport := &http.Transport{Proxy: http.ProxyURL(n.ProxURL)}
	client := http.Client{Transport: transport}

	// client := http.Client{}
	request, err := n.GetRequest()
	if err != nil {
		err = fmt.Errorf("构造请求报文错误: %w", err)
		n.CurRequestMessage = err.Error() // 更改请求报文为错误信息
		n.CurResponseMessage = "无响应报文"    // 更改响应报文为错误信息
		return ResponseJSON{StatusCode: 500, Message: err.Error()}
	}
	// 解析请求信息
	RequestMessage, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		n.CurRequestMessage = "解析请求报文错误" + err.Error() // 更改请求报文为错误信息
		n.CurResponseMessage = "无响应报文"                 // 更改响应报文为错误信息
	} else {
		n.CurRequestMessage = string(RequestMessage)
	}

	resp, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("发送请求错误: %w", err)
		n.CurResponseMessage = err.Error()
		return ResponseJSON{StatusCode: 500, Message: err.Error()}
	}
	// 解析响应信息
	ResponseMessage, err := httputil.DumpResponse(resp, true)
	if err != nil {
		n.CurResponseMessage = "解析响应包错误" + err.Error()
	} else {
		n.CurResponseMessage = string(ResponseMessage)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("解析响应体错误: %w", err)
		n.CurResponseMessage = "解析响应体错误" + err.Error()
		return ResponseJSON{StatusCode: 500, Message: err.Error()}
	}

	n.ResponsePool[fmt.Sprint(n.Interface.Serial)] = string(respBytes)
	n.ResponsePoolType[fmt.Sprint(n.Interface.Serial)] = n.Interface.ResponseType

	// 如果接口响应需要校验
	if n.Interface.ResponseValidate {
		retVal, err := n.GetResponseMappedValue()
		if err != nil {
			err = fmt.Errorf("校验响应错误: %w", err)
			n.CurResponseMessage = err.Error() + "\n响应包如下\n\n" + n.CurResponseMessage
			return ResponseJSON{StatusCode: 500, Message: err.Error()}
		}
		if retVal == "" {
			err = fmt.Errorf("校验响应为空: %w", err)
			n.CurResponseMessage = err.Error() + "\n响应包如下\n\n" + n.CurResponseMessage
			return ResponseJSON{StatusCode: 500, Message: err.Error()}
		}
	}

	return ResponseJSON{StatusCode: 200, Message: string(respBytes), Data: string(respBytes)}

}
