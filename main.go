package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// Version is populated at compile time via -ldflags "-X main.Version=x.y.z" in CI
var Version = "dev"

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:         "SimpleAC",
		Width:         450,
		Height:        600,
		MinWidth:      450,
		MinHeight:     600,
		MaxWidth:      450,
		MaxHeight:     600,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
