package main

import (
	"os"

	"gopkg.in/ini.v1"
)

const settingsFile = "settings.ini"

// LoadSettings reads the plugin statuses from the settings.ini file
// If the file does not exist, it creates an empty one with all plugins enabled by default
func LoadSettings(filename string) (map[string]bool, error) {
	pluginStatus := make(map[string]bool)

	// Check if the settings file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// If the file does not exist, initialize all plugins to true by default
		for _, plugin := range GetPlugins() {
			pluginStatus[plugin.Name] = true
		}
		// Create and save the default settings file
		err = SaveSettings(filename, pluginStatus)
		if err != nil {
			return nil, err
		}
		return pluginStatus, nil
	}

	// Load the settings file
	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}

	pluginsSection := cfg.Section("plugins")
	for _, key := range pluginsSection.Keys() {
		pluginStatus[key.Name()] = key.MustBool(true)
	}
	return pluginStatus, nil
}

// SaveSettings writes the plugin statuses to the settings.ini file
func SaveSettings(filename string, pluginStatus map[string]bool) error {
	cfg := ini.Empty()
	section, err := cfg.NewSection("plugins")
	if err != nil {
		return err
	}
	for plugin, status := range pluginStatus {
		// prettyName := GetPrettyName(plugin)
		_, err := section.NewKey(plugin, boolToStr(status))
		if err != nil {
			return err
		}
	}
	return cfg.SaveTo(filename)
}

// Helper function to convert bool to string
func boolToStr(val bool) string {
	if val {
		return "true"
	}
	return "false"
}
