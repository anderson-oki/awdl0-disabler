package ui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Header        lipgloss.Style
	Footer        lipgloss.Style
	StatusUp      lipgloss.Style
	StatusDown    lipgloss.Style
	StatusUnknown lipgloss.Style
	Bar           lipgloss.Style
	Dashboard     lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1),
		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(1, 0),
		StatusUp: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true),
		StatusDown: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
		StatusUnknown: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Bold(true),
		Bar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")),
		Dashboard: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(2, 4),
	}
}
