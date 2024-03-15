package utils

import (
	"fmt"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// RODController 机器人控制器
type RODController struct {
	Browser   *rod.Browser       // 浏览器
	Launcher  *launcher.Launcher // 启动器
	CheckTime int                // 检查时间
}

// StartBrowser 开启一个浏览器
func (rodc *RODController) StartBrowser() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// 设置浏览器启动器
		browserPath, ok := launcher.LookPath() // 查找浏览器路径
		fmt.Println("启用本地浏览器: ", ok)
		if ok { // 如果找到本地浏览器，则使用
			rodc.Launcher = launcher.New().Bin(browserPath) // 设置浏览器
		} else { // 如果未找到则不使用本地浏览器
			rodc.Launcher = launcher.New() // 设置浏览器
		}
		rodc.Launcher.Headless(false) // 设置浏览器参数
		defer rodc.Launcher.Cleanup() // cleanup

		// 启动浏览器
		rodc.Browser = rod.New().NoDefaultDevice().ControlURL(rodc.Launcher.MustLaunch()).MustConnect() // 打开浏览器
		fmt.Println("浏览器已启动")
		defer wg.Done() // 判断是否执行完毕
	}()
	wg.Wait()
}

// CloseBrowser 关闭浏览器
func (rodc *RODController) CloseBrowser() {
	rodc.Launcher.Kill()
	rodc.Browser = nil
	rodc.Launcher = nil
}
