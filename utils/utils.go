package utils

import "github.com/charmbracelet/bubbletea"

// Wrap casts a message into a tea.Cmd
func Wrap(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
