package main

import (
	"github.com/getlantern/systray"
)

var isEnabled = true

func OnReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Clipboard Master")
	systray.SetTooltip("Clipboard Master")

	mEnable := systray.AddMenuItemCheckbox(
		getLocalization("systray.enable.label"),
		getLocalization("systray.enable.tooltip"),
		isEnabled,
	)
	mStartup := systray.AddMenuItemCheckbox(
		getLocalization("systray.bootAtStartup.label"),
		getLocalization("systray.bootAtStartup.tooltip"),
		IsStartupEnabled(),
	)
	mSettings := systray.AddMenuItem(
		getLocalization("systray.settings.label"),
		getLocalization("systray.settings.tooltip"),
	)
	mPluginSettings := make(map[string]*systray.MenuItem)

	for _, plugin := range GetPlugins() {
		prettyName := GetPrettyName(plugin.Name)
		enabled := pluginStatus[plugin.Name]
		mPluginSettings[plugin.Name] = mSettings.AddSubMenuItemCheckbox(
			prettyName,
			getLocalization("systray.settingsToggle.label")+prettyName,
			enabled,
		)
	}

	mQuit := systray.AddMenuItem(
		getLocalization("systray.quit.label"),
		getLocalization("systray.quit.tooltip"),
	)

	go func() {
		MonitorClipboard(GetPlugins(), pluginStatus)
	}()

	go func() {
		for {
			<-mStartup.ClickedCh
			if mStartup.Checked() {
				err := SetStartup(false)
				if err != nil {
					SendNotification(
						getLocalization("notifications.errorSettingStartup.title"),
						getLocalization("notifications.errorSettingStartup.message")+err.Error(),
					)
				} else {
					mStartup.Uncheck()
				}
			} else {
				err := SetStartup(true)
				if err != nil {
					SendNotification(
						getLocalization("notifications.errorSettingStartup.title"),
						getLocalization("notifications.errorSettingStartup.message")+err.Error(),
					)
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
					err := SavePluginSettings(pluginStatus)
					if err != nil {
						SendNotification(
							getLocalization("notifications.errorSavingSettings.title"),
							getLocalization("notifications.errorSavingSettings.message")+err.Error(),
						)
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
	err := SavePluginSettings(pluginStatus)
	if err != nil {
		SendNotification(
			getLocalization("notifications.errorSavingSettings.title"),
			getLocalization("notifications.errorSavingSettings.message")+err.Error(),
		)
	}
}
