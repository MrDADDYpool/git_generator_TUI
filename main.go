package main

import (
	"fmt"
	"os"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	model.go
	update.go
	view.go
)

func main() {
	// Check the operating system
	if runtime.GOOS == "windows" {
		fmt.Println("Note: This script is running on Windows. Some commands might differ.")
	}

	p := tea.NewProgram(model.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
