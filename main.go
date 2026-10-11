package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
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
		// Matches --bg-window (light) so there is no coloured flash before the UI paints.
		BackgroundColour: &options.RGBA{R: 239, G: 239, B: 241, A: 255},
		Mac: &mac.Options{
			// Traffic lights sit inside the app's own 52px toolbar, like native macOS apps.
			TitleBar: mac.TitleBarHiddenInset(),
			// Follow the system Light/Dark appearance; the UI can still override it in Settings.
			Appearance: mac.DefaultAppearance,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
