package main

import (
	_ "embed"
)

// Embed the icon file
//
//go:embed assets/icon.ico
var iconData []byte
