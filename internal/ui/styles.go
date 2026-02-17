package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("62")
	ColorSuccess   = lipgloss.Color("42")
	ColorWarning   = lipgloss.Color("214")
	ColorError     = lipgloss.Color("196")
	ColorDim       = lipgloss.Color("59")
	ColorBg        = lipgloss.Color("235")

	// Styles
	StyleTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(ColorPrimary).
		Padding(0, 2)

	StyleSelected = lipgloss.NewStyle().
		Background(ColorPrimary).
		Foreground(lipgloss.Color("230"))

	StyleDim = lipgloss.NewStyle().
		Foreground(ColorDim)

	StyleSuccess = lipgloss.NewStyle().
		Foreground(ColorSuccess)

	StyleError = lipgloss.NewStyle().
		Foreground(ColorError)

	StyleWarning = lipgloss.NewStyle().
		Foreground(ColorWarning)

	StyleCompleted = lipgloss.NewStyle().
		Foreground(ColorDim).
		Strikethrough(true)

	StyleHighPriority = lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true)

	StyleMediumPriority = lipgloss.NewStyle().
		Foreground(ColorWarning)

	StyleTag = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("5")).
		Padding(0, 1)
)