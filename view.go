package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	model.go
)

var (
	titleASCII = `
  ________._________________ ___  ____ _____________     _____      _____    _______      _____    _____________________________ 
 /  _____/|   \__    ___/   |   \|    |   \______   \   /     \    /  _  \   \      \    /  _  \  /  _____/\_   _____/\______   \
/   \  ___|   | |    | /    ~    \    |   /|    |  _/  /  \ /  \  /  /_\  \  /   |   \  /  /_\  \/   \  ___ |    __)_  |       _/
\    \_\  \   | |    | \    Y    /    |  / |    |   \ /    Y    \/    |    \/    |    \/    |    \    \_\  \|        \ |    |   \
 \______  /___| |____|  \___|_  /|______/  |______  / \____|__  /\____|__  /\____|__  /\____|__  /\______  /_______  / |____|_  /
        \/                    \/                  \/          \/         \/         \/         \/        \/        \/         \/ 
`

	optionStyle = lipgloss.NewStyle().PaddingLeft(2)

	selectedOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Background(lipgloss.Color("57")).
				PaddingLeft(2).
				Bold(true)

	bracketStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).SetString("X")
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func View(m model.Model) string {
	var b strings.Builder

	b.WriteString(titleASCII + "\n\n")

	if m.InputMode {
		b.WriteString(inputStyle.Render(fmt.Sprintf("%s: %s", m.InputField, m.InputValue)) + "\n\n")
	}

	for i, option := range m.Options {
		if m.CurrentOption == i {
			b.WriteString(bracketStyle.Render("[") + cursorStyle.String() + bracketStyle.Render("] ") + selectedOptionStyle.Render(option) + "\n")
		} else {
			b.WriteString(bracketStyle.Render("[ ] ") + optionStyle.Render(option) + "\n")
		}
	}

	return b.String()
}
