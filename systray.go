package main

import (
	"github.com/getlantern/systray"
)

func OnReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Clipboard Master")
	systray.SetTooltip("Clipboard Master")

	mEnable := systray.AddMenuItemCheckbox("Enable", "Enable/Disable the app", isEnabled)
	mStartup := systray.AddMenuItemCheckbox("Boot at Startup", "Enable/Disable boot at startup", IsStartupEnabled())
	mSettings := systray.AddMenuItem("Settings", "Settings")
	mPluginSettings := make(map[string]*systray.MenuItem)

	for _, plugin := range GetPlugins() {
		prettyName := GetPrettyName(plugin.Name)
		enabled := pluginStatus[plugin.Name]
		mPluginSettings[plugin.Name] = mSettings.AddSubMenuItemCheckbox(prettyName, "Enable/Disable "+prettyName, enabled)
	}

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		MonitorClipboard(GetPlugins(), pluginStatus)
	}()

	go func() {
		for {
			<-mStartup.ClickedCh
			if mStartup.Checked() {
				err := SetStartup(false)
				if err != nil {
					SendNotification("Error", "Error setting startup: %v"+err.Error())
				} else {
					mStartup.Uncheck()
				}
			} else {
				err := SetStartup(true)
				if err != nil {
					SendNotification("Error", "Error setting startup: %v"+err.Error())
				} else {
					mStartup.Check()
				}
			}
		}
	}()

	go func() {
		for {
			<-mEnable.ClickedCh
			if mEnable.Checked() {
				isEnabled = false
				mEnable.Uncheck()
			} else {
				isEnabled = true
				mEnable.Check()
			}
		}
	}()

	go func() {
		for pluginName, menuItem := range mPluginSettings {
			go func(name string, item *systray.MenuItem) {
				for {
					<-item.ClickedCh
					pluginStatus[name] = !pluginStatus[name]
					if item.Checked() {
						item.Uncheck()
					} else {
						item.Check()
					}
					err := SaveSettings(settingsFile, pluginStatus)
					if err != nil {
						SendNotification("Error", "Error saving settings: %v"+err.Error())
					}
				}
			}(pluginName, menuItem)
		}
	}()

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func OnExit() {
	err := SaveSettings(settingsFile, pluginStatus)
	if err != nil {
		SendNotification("Error", "Error saving settings: %v"+err.Error())
	}
}
