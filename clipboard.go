package main

import (
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func MonitorClipboard(plugins []Plugin, pluginStatus map[string]bool) {
	var processedText string
	for {
		if !isEnabled {
			time.Sleep(1 * time.Second)
			continue
		}

		text, err := clipboard.ReadAll()

		if err != nil || !IsValidURL(text) {
			time.Sleep(1 * time.Second)
			continue
		}

		if text != processedText {
			if ContainsReplacement(text) {
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
				if strings.Contains(processedText, "?") {
					processedText = processedText[:strings.Index(processedText, "?")]
				}
				clipboard.WriteAll(processedText)
				SendNotification(
					getLocalization("notifications.clipboardUpdated.title"),
					processedText,
				)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
