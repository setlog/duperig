package main

import (
	"fmt"

	"github.com/fatih/color"
)

func colorize(text string, colorCode color.Attribute, clearAfter bool) string {
	if text == "" {
		if clearAfter {
			return clearStyle()
		}
		return text
	}
	if clearAfter {
		return fmt.Sprintf("%s%s%s", formatStyle(colorCode), text, clearStyle())
	}
	return fmt.Sprintf("%s%s", formatStyle(colorCode), text)
}

func formatStyle(colorCode color.Attribute) string {
	return fmt.Sprintf("\x1B[%dm", colorCode)
}

func clearStyle() string {
	return fmt.Sprintf("\x1B[%dm", color.Reset)
}
