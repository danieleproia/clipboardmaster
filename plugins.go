package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

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

var (
	plugins      []Plugin
	prettyToNorm map[string]string
	normToPretty map[string]string
)

// LoadPlugins loads all plugins from the specified directory
func LoadPlugins(dir string) error {
	SendNotification("Updating plugins", "Reading plugins...")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	prettyToNorm = make(map[string]string)
	normToPretty = make(map[string]string)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml" {
			data, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return err
			}

			var plugin Plugin
			err = yaml.Unmarshal(data, &plugin)
			if err != nil {
				return err
			}

			// Normalize the plugin name
			normName := normalizeName(plugin.Name)
			prettyToNorm[plugin.Name] = normName
			normToPretty[normName] = plugin.Name

			plugin.Name = normName
			plugins = append(plugins, plugin)
		}
	}
	SendNotification("Updated list of plugins", fmt.Sprintf("Found %d plugins", len(plugins)))
	return nil
}

// GetPlugins returns the list of loaded plugins
func GetPlugins() []Plugin {
	return plugins
}

// ContainsReplacement checks if the text contains any replacement strings
func ContainsReplacement(text string) bool {
	for _, plugin := range plugins {
		for _, transformation := range plugin.Transformations {
			if strings.Contains(text, transformation.Replace) {
				return true
			}
		}
	}
	return false
}

// normalizeName converts a string to lowercase and replaces spaces with underscores
func normalizeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

// GetPrettyName returns the pretty name for a normalized name
func GetPrettyName(normName string) string {
	return normToPretty[normName]
}

// GetNormName returns the normalized name for a pretty name
func GetNormName(prettyName string) string {
	return prettyToNorm[prettyName]
}
