package main

import (
	"os"

	"gopkg.in/ini.v1"
)

const settingsFile = "settings.ini"

func createSettingsFile() error {
	// if file empty, create it
	var cfg *ini.File = ini.Empty()
	var err error
	var languageSection *ini.Section
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		// Create the settings file
		file, err := os.Create(settingsFile)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		// Load the settings file
		cfg, err = ini.Load(settingsFile)
		if err != nil {
			return err
		}
	}
	// Check if the plugins section exists
	_, err = cfg.GetSection("plugins")
	if err != nil {
		_, err = cfg.NewSection("plugins")
		if err != nil {
			return err
		}
	}
	for _, plugin := range GetPlugins() {
		_, err = cfg.Section("plugins").NewKey(plugin.Name, "true")
		if err != nil {
			return err
		}
	}

	languageSection, err = cfg.NewSection("language")
	if err != nil {
		return err
	}
	_, err = languageSection.GetKey("language")
	if err != nil {
		_, err = languageSection.NewKey("language", "en")
		if err != nil {
			return err
		}
	}
	cfg.SaveTo(settingsFile)
	return nil
}

// LoadSettings reads the plugin statuses from the settings.ini file
// If the file does not exist, it creates an empty one with all plugins enabled by default
func LoadSettings() (map[string]bool, error) {
	createSettingsFile()
	pluginStatus := make(map[string]bool)

	// Load the settings file
	cfg, err := ini.Load(settingsFile)
	if err != nil {
		return nil, err
	}

	pluginsSection := cfg.Section("plugins")
	// if the plugins section has no keys, create them
	if len(pluginsSection.Keys()) == 0 {
		for _, plugin := range GetPlugins() {
			_, err := pluginsSection.NewKey(plugin.Name, "true")
			if err != nil {
				return nil, err
			}
		}
	}

	for _, key := range pluginsSection.Keys() {
		pluginStatus[key.Name()] = key.MustBool(true)
	}

	languageSection := cfg.Section("language")
	language = languageSection.Key("language").String()

	return pluginStatus, nil
}

// SaveSettings writes the plugin statuses to the settings.ini file
func SavePluginSettings(pluginStatus map[string]bool) error {
	cfg, err := ini.Load(settingsFile)
	if err != nil {
		return err
	}
	section, err := cfg.GetSection("plugins")
	if err != nil {
		return err
	}
	for plugin, status := range pluginStatus {
		if !section.HasKey(plugin) {
			_, err := section.NewKey(plugin, boolToStr(status))
			if err != nil {
				return err
			}
		} else {
			section.Key(plugin).SetValue(boolToStr(status))
		}
	}
	return cfg.SaveTo(settingsFile)
}

// Helper function to convert bool to string
func boolToStr(val bool) string {
	if val {
		return "true"
	}
	return "false"
}
