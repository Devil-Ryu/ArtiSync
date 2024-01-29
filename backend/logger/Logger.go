package artlog

import (
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 日志等级枚举
const (
	INFO  string = "INF" // 一般信息
	DEBUG string = "DEB" // 调试信息
	WARN  string = "WAR" // 警告信息
	ERROR string = "ERR" // 错误信息
)

// Logger 日志
type Logger struct {
	Ctx     context.Context
	History []LoggerResult // 日志历史
	Level   string         // 日志等级
	Push    bool           // 是否推送给前端
}

// LoggerResult 日志返回
type LoggerResult struct {
	Datetime string `json:"datetime"`
	Level    string `json:"level"`
	Tag      string `json:"tag"`
	Message  string `json:"message"`
}

// NewLogger 日志构造方法
func NewLogger() *Logger {
	return &Logger{}
}

// Startup wails设置context
func (l *Logger) Startup(ctx context.Context) {
	l.Ctx = ctx
	l.Push = true
}

// Print 显示日志
func (l *Logger) Print(level string, tag string, message string) {
	curTime := time.Now().Format("2006-01-02 15:04:05")                              // 获取当前时间
	log := LoggerResult{Datetime: curTime, Level: level, Tag: tag, Message: message} // 生成日志
	fmt.Printf("%s | %s | %s | %s\n", log.Datetime, log.Level, log.Tag, log.Message) // 格式化打印日志
	l.History = append(l.History, log)

	// 推送日志
	if l.Push {
		l.PushLog(log)
	}
}

// PushLog 推送日志的方法，可以自定义，此处默认推送给前端PushLog方法
func (l *Logger) PushLog(log LoggerResult) {
	runtime.EventsEmit(l.Ctx, "PushLog", log)
}
