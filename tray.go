package main

import (
	"context"
	_ "embed"
	"os"
	"time"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/appicon.png
var trayIcon []byte

func (a *App) InitSystray(ctx context.Context) {
	systray.Run(func() {
		systray.SetIcon(trayIcon)
		systray.SetTooltip("SimpleAC VPN")

		show := systray.AddMenuItem("Show", "Show The Window")
		show.Click(func() { runtime.WindowShow(ctx) })

		systray.AddSeparator()

		mStatus := systray.AddMenuItem("⚪ Status: Unknown", "Current connection status")
		mStatus.Disable()

		mToggle := systray.AddMenuItem("Connect", "Toggle VPN Connection")

		mToggle.Click(func() {
			// Must run in a goroutine — tray callbacks execute on the macOS
			// main thread, and Connect/Disconnect are blocking subprocess calls.
			// Calling them directly here deadlocks / panics the app.
			go func() {
				connected, _ := a.IsConnected()
				if connected {
					_, _ = a.Disconnect()
				} else {
					_, _ = a.Connect()
				}
			}()
		})

		systray.AddSeparator()

		systray.SetOnClick(func(menu systray.IMenu) { menu.ShowMenu() })

		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		mQuit.Click(func() {
			systray.Quit()
			os.Exit(0)
		})

		// Background goroutine to check connection status
		// Using global connection state here
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				a.mu.RLock()
				hasClient := a.VPNClient != nil
				a.mu.RUnlock()

				if !hasClient {
					mStatus.SetTitle("🔴 Not Configured")
					mToggle.Disable()
					mToggle.SetTitle("Connect")
					continue
				}

				mToggle.Enable()
				connected, _ := a.IsConnected()
				isProcessing := a.IsConnecting()

				if isProcessing {
					if connected {
						mStatus.SetTitle("🟠 Disconnecting...")
						mToggle.SetTitle("Disconnecting...")
					} else {
						mStatus.SetTitle("🟠 Connecting...")
						mToggle.SetTitle("Connecting...")
					}
					mToggle.Disable()
				} else {
					if connected {
						mStatus.SetTitle("🟢 Connected")
						mToggle.SetTitle("Disconnect")
					} else {
						mStatus.SetTitle("🔴 Disconnected")
						mToggle.SetTitle("Connect")
					}
				}
			}
		}()

	}, func() {
		// On exit cleanup
	})
}
