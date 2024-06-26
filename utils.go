package main

import (
	_ "embed"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

//go:embed assets/icon.ico
var iconData []byte

func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	// A valid URL should have a scheme and a host
	return u.Scheme != "" && u.Host != ""
}

func SetStartup(enable bool) error {
	key := `Software\Microsoft\Windows\CurrentVersion\Run`

	exePath, err := os.Executable()
	if err != nil {
		SendNotification(
			getLocalization("notifications.errorSettingStartup.title"),
			getLocalization("notifications.errorSettingStartup.message")+err.Error(),
		)
		return err
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		SendNotification(
			getLocalization("notifications.errorSettingStartup.title"),
			getLocalization("notifications.errorSettingStartup.message")+err.Error(),
		)
		return err
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, key, registry.SET_VALUE)
	if err != nil {
		SendNotification(
			getLocalization("notifications.errorSettingStartup.title"),
			getLocalization("notifications.errorSettingStartup.message")+err.Error(),
		)
		return err
	}
	defer k.Close()

	if enable {
		err = k.SetStringValue(appName, exePath)
		if err != nil {
			SendNotification(
				getLocalization("notifications.errorSettingStartup.title"),
				getLocalization("notifications.errorSettingStartup.message")+err.Error(),
			)
			return err
		}
	} else {
		err = k.DeleteValue(appName)
		if err != nil {
			SendNotification(
				getLocalization("notifications.errorSettingStartup.title"),
				getLocalization("notifications.errorSettingStartup.message")+err.Error(),
			)
			return err
		}
	}
	return nil
}

func IsStartupEnabled() bool {
	key := `Software\Microsoft\Windows\CurrentVersion\Run`

	k, err := registry.OpenKey(registry.CURRENT_USER, key, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	_, _, err = k.GetStringValue(appName)
	return err == nil
}
