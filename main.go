package main

import (
	"os"

	"github.com/getlantern/systray"
)

const appName = "ClipboardMaster"

var exePath, _ = os.Executable()

var isEnabled = true
var pluginStatus map[string]bool

func main() {
	var err error
	// Load plugins from the plugins directory
	err = LoadPlugins("./plugins")
	if err != nil {
		SendNotification("Error", "Error loading plugins: %v"+err.Error())
	}

	// Load settings from the settings.ini file
	pluginStatus, err = LoadSettings(settingsFile)
	if err != nil {
		SendNotification("Error", "Error loading settings: %v"+err.Error())
	}

	// Initialize plugin status map if settings not loaded
	if pluginStatus == nil {
		pluginStatus = make(map[string]bool)
		for _, plugin := range GetPlugins() {
			pluginStatus[plugin.Name] = true
		}
	}

	// Start the system tray
	systray.Run(OnReady, OnExit)
}
