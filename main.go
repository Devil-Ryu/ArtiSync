package main

import (
	"ArtiSync/backend/api"
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	dbController := api.NewDBController()
	atController := api.NewATController()
	atController.SetDBController(dbController) // SetDBController设置数据库控制器
	atController.InitConfig()                  // 初始化文章控制器的配置

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ArtiSync",
		Width:  1224,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			atController.Startup(ctx)
			dbController.Startup(ctx)
		},
		Bind: []interface{}{
			dbController,
			atController,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
