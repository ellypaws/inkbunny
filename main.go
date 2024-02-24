package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"inkbunny/gui"
	"log"
	"os"
)

func main() {
	if _, err := tea.NewProgram(
		gui.InitialModel(os.Getenv("SID")),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
