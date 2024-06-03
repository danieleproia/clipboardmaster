package main

import (
	"log"
	"path/filepath"

	"github.com/go-toast/toast"
)

func SendNotification(title string, message string) {
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
