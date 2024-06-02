package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"github.com/go-toast/toast"
	"golang.org/x/sys/windows/registry"
	"gopkg.in/yaml.v3"
)

const appName = "ClipboardMaster"

var exePath, _ = os.Executable()

var isEnabled = true

// Transformation represents a single find-replace operation
type Transformation struct {
	Find    string `yaml:"find"`
	Replace string `yaml:"replace"`
}

// Plugin represents a plugin described in a YAML file
type Plugin struct {
	Name            string           `yaml:"name"`
	Transformations []Transformation `yaml:"replacements"`
}

func isValidURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	// A valid URL should have a scheme and a host
	return u.Scheme != "" && u.Host != ""
}

// LoadPlugins loads all plugins from the specified directory
func LoadPlugins(dir string) ([]Plugin, error) {
	var plugins []Plugin

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			data, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}

			var plugin Plugin
			err = yaml.Unmarshal(data, &plugin)
			if err != nil {
				return nil, err
			}

			plugins = append(plugins, plugin)
		}
	}

	return plugins, nil
}

// Check if the text contains any of the replacement strings
func containsReplacement(text string, plugins []Plugin) bool {
	for _, plugin := range plugins {
		for _, transformation := range plugin.Transformations {
			if strings.Contains(text, transformation.Replace) {
				return true
			}
		}
	}
	return false
}

func sendNotification(message string) {
	iconPath := filepath.Dir(exePath)
	iconPath = iconPath + "\\icon.png"
	fmt.Println(iconPath)
	notification := toast.Notification{
		AppID:   "Clipboard Master",
		Title:   "Clipboard Updated",
		Message: message,
		Icon:    iconPath,
	}

	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

func monitorClipboard(plugins []Plugin) {
	var previousText string
	var processedText string

	for {
		// if app is disabled, sleep for a while and check again
		if !isEnabled {
			time.Sleep(1 * time.Second)
			continue
		}
		// Get the current clipboard content
		text, err := clipboard.ReadAll()
		if err != nil || !isValidURL(text) {
			time.Sleep(1 * time.Second)
			continue
		}

		// If the clipboard content has changed and hasn't been processed already
		if text != previousText && text != processedText {
			// Check if the text contains any replacement strings
			if containsReplacement(text, plugins) {
				previousText = text
				continue
			}

			processedText = text
			for _, plugin := range plugins {
				for _, transformation := range plugin.Transformations {
					processedText = strings.ReplaceAll(processedText, transformation.Find, transformation.Replace)
				}
			}
			if processedText != text {
				// remove everything after the ? in the url
				if strings.Contains(processedText, "?") {
					processedText = processedText[:strings.Index(processedText, "?")]
				}
				// Update the clipboard with the processed text
				clipboard.WriteAll(processedText)
				// Send a notification
				sendNotification(processedText)
			}
			previousText = text
		}

		// Sleep for a short duration before checking again
		time.Sleep(500 * time.Millisecond)
	}
}

func setStartup(enable bool) error {
	key := `Software\Microsoft\Windows\CurrentVersion\Run`

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to get executable path: %v", err)
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("unable to get absolute path of executable: %v", err)
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, key, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("unable to open registry key: %v", err)
	}
	defer k.Close()

	if enable {
		err = k.SetStringValue(appName, exePath)
		if err != nil {
			return fmt.Errorf("unable to set registry value: %v", err)
		}
	} else {
		err = k.DeleteValue(appName)
		if err != nil {
			return fmt.Errorf("unable to delete registry value: %v", err)
		}
	}

	return nil
}

func isStartupEnabled() bool {
	key := `Software\Microsoft\Windows\CurrentVersion\Run`

	k, err := registry.OpenKey(registry.CURRENT_USER, key, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	_, _, err = k.GetStringValue(appName)
	return err == nil
}

func main() {
	// Load plugins from the plugins directory
	plugins, err := LoadPlugins("./plugins")
	if err != nil {
		sendNotification("Error loading plugins: %v" + err.Error())
	}

	// Start the system tray
	systray.Run(onReady(plugins), onExit)
}

func onReady(plugins []Plugin) func() {
	return func() {
		// Set the icon
		systray.SetIcon(iconData)
		systray.SetTitle("Clipboard Master")
		systray.SetTooltip("Clipboard Master")

		// add a checkbox to enable/disable the app
		mEnable := systray.AddMenuItemCheckbox("Enable", "Enable/Disable the app", isEnabled)
		// add a checkbox for boot at startup
		mStartup := systray.AddMenuItemCheckbox("Boot at Startup", "Enable/Disable boot at startup", isStartupEnabled())
		// Add a quit button
		mQuit := systray.AddMenuItem("Quit", "Quit the application")

		// Start monitoring the clipboard in a separate goroutine
		go func() {
			monitorClipboard(plugins)
		}()

		// Handle startup checkbox click
		go func() {
			for {
				<-mStartup.ClickedCh
				if mStartup.Checked() {
					err := setStartup(false)
					if err != nil {
						sendNotification("Error setting startup: %v" + err.Error())
					} else {
						mStartup.Uncheck()
					}
				} else {
					err := setStartup(true)
					if err != nil {
						sendNotification("Error setting startup: %v" + err.Error())
					} else {
						mStartup.Check()
					}
				}
			}
		}()

		// Handle enable checkbox click
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

		// Handle quit button click
		go func() {
			<-mQuit.ClickedCh
			systray.Quit()
			fmt.Println("Quit")
		}()
	}
}

func onExit() {
	// Clean up here if needed
}
