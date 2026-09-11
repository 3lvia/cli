package style

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

type PrintOptions struct {
	Color string
}

func Print(message string, options *PrintOptions) {
	if options == nil {
		options = &PrintOptions{}
	}

	colorCode := getColorCode(options.Color)
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorCode))

	fmt.Println(style.Render(message))
}

func PrintInfo(message string) {
	Print(message, &PrintOptions{Color: "blue"})
}

func PrintSuccess(message string) {
	Print(message, &PrintOptions{Color: "green"})
}

func PrintWarning(message string) {
	Print(message, &PrintOptions{Color: "yellow"})
}

func PrintError(message string) {
	Print(message, &PrintOptions{Color: "red"})
}

func getColorCode(color string) string {
	switch color {
	case "red":
		return "1"
	case "green":
		return "2"
	case "yellow":
		return "3"
	case "blue":
		return "4"
	case "magenta":
		return "5"
	case "cyan":
		return "6"
	case "white":
		return "7"
	default:
		return "7"
	}
}
