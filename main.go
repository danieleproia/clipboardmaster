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
var plugins []Plugin
var pluginStatus map[string]bool

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

	// notification of reading plugins
	sendNotification("Updating plugins", "Reading plugins...")
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
	// notification found n plugins
	sendNotification("Updated list of plugins", fmt.Sprintf("Found %d plugins", len(plugins)))
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

func sendNotification(title string, message string) {
	duro, _ := toast.Duration("short")
	iconPath := filepath.Dir(exePath)
	iconPath = iconPath + "\\icon.png"
	notification := toast.Notification{
		AppID:    "Clipboard Master",
		Title:    title,
		Message:  message,
		Icon:     iconPath,
		Duration: duro,
	}

	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

func monitorClipboard(plugins []Plugin, pluginStatus map[string]bool) {
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
		if text != processedText {
			// Check if the text contains any replacement strings
			if containsReplacement(text, plugins) {
				continue
			}
			processedText = text
			for _, plugin := range plugins {
				if !pluginStatus[plugin.Name] {
					continue
				}
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
				sendNotification("Updated clipboard", processedText)
			}
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
	var err error
	// Load plugins from the plugins directory
	plugins, err = LoadPlugins("./plugins")
	if err != nil {
		sendNotification("Error", "Error loading plugins: %v"+err.Error())
	}
	// Initialize plugin status map
	pluginStatus = make(map[string]bool)
	for _, plugin := range plugins {
		pluginStatus[plugin.Name] = true
	}

	// Start the system tray
	systray.Run(onReady, onExit)
}

func onReady() {
	// Set the icon
	systray.SetIcon(iconData)
	systray.SetTitle("Clipboard Master")
	systray.SetTooltip("Clipboard Master")

	// add a checkbox to enable/disable the app
	mEnable := systray.AddMenuItemCheckbox("Enable", "Enable/Disable the app", isEnabled)
	// add a checkbox for boot at startup
	mStartup := systray.AddMenuItemCheckbox("Boot at Startup", "Enable/Disable boot at startup", isStartupEnabled())
	// Add settings menu for plugins
	mSettings := systray.AddMenuItem("Settings", "Settings")
	mPluginSettings := make(map[string]*systray.MenuItem)

	for _, plugin := range plugins {
		mPluginSettings[plugin.Name] = mSettings.AddSubMenuItemCheckbox(plugin.Name, "Enable/Disable "+plugin.Name, true)
	}

	// Add a quit button
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Start monitoring the clipboard in a separate goroutine
	go func() {
		monitorClipboard(plugins, pluginStatus)
	}()

	// Handle startup checkbox click
	go func() {
		for {
			<-mStartup.ClickedCh
			if mStartup.Checked() {
				err := setStartup(false)
				if err != nil {
					sendNotification("Error", "Error setting startup: %v"+err.Error())
				} else {
					mStartup.Uncheck()
				}
			} else {
				err := setStartup(true)
				if err != nil {
					sendNotification("Error", "Error setting startup: %v"+err.Error())
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

	// Handle plugin settings checkbox click
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
				}
			}(pluginName, menuItem)
		}
	}()

	// Handle quit button click
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// Clean up here if needed
}
