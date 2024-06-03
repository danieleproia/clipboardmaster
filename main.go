package main

import (
	"os"

	"github.com/getlantern/systray"
)

const appName = "ClipboardMaster"

var exePath, _ = os.Executable()

func main() {
	var err error
	// Load settings from the settings.ini file
	pluginStatus, err = LoadSettings()
	if err != nil {
		SendNotification(
			getLocalization("notifications.errorLoadingSettings.title"),
			getLocalization("notifications.errorLoadingSettings.message")+err.Error(),
		)
	}
	// Load localization from the en.po file
	lang = generateLocalization()

	// Load plugins from the plugins directory
	err = LoadPlugins()
	if err != nil {
		SendNotification(
			getLocalization("notifications.errorLoadingPlugins.title"),
			getLocalization("notifications.errorLoadingPlugins.message")+err.Error(),
		)
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
