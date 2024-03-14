package api

import (
	artlog "ArtiSync/backend/logger"
)

// RODController 机器人控制器
type RODController struct {
	ArticleList []Article // 文章列表
	// ProxyURL          *url.URL                 // 网络代理
	RequestSleep  int           // 请求一次休眠时间
	ImageReadType string        // 图片读取方式
	DBController  *DBController // 数据库控制器
	artlog.Logger               // 检查时间
}
