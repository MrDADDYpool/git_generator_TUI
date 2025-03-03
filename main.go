package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

// General stuff for styling the view
var (
	keywordStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	mainStyle           = lipgloss.NewStyle().MarginLeft(2)
	optionStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Background(lipgloss.Color("57")).
				PaddingLeft(2).
				Bold(true)
	bracketStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).SetString("X")
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type model struct {
	Choice         int
	ShowInputField bool
	InputField     textinput.Model
}

const (
	optionCreateSSHKey       = "Create SSH Key"
	optionSetGlobalGitConfig = "Set Global Git Config"
	optionCloneGitHubRepo    = "Clone GitHub Repository"
	optionCommitAndSync      = "Commit and Sync Changes"
	optionExit               = "Exit"
	optionBack               = "Back"
	optionCancel             = "Cancel"
	optionFilePath           = "Enter file path"
	optionPassphrase         = "Enter passphrase"
	optionGenerateKeys       = "Generate Keys"
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
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "https://github.com/lol"
	ti.CharLimit = 156
	ti.Width = 20
	ti.Cursor.Style = cursorStyle

	return model{
		Choice:         0,
		ShowInputField: false,
		InputField:     ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// A chaque event, qu'est-ce qu'on fait sur le model ?
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp:
			if m.Choice > 0 {
				m.Choice--
			} else {
				m.Choice = 2
			}
		case tea.KeyDown:
			if m.Choice >= 2 {
				m.Choice = 0
			} else {
				m.Choice++
			}
		case tea.KeyEnter:
			if m.ShowInputField {
				m.ShowInputField = false
			} else {
				m.ShowInputField = true
				m.InputField.Focus()
			}

			return m, performAction(m.Choice)
		}

	// We handle errors just like any other message
	case errMsg:
		return m, tea.Quit
	}

	if m.ShowInputField {
		m.InputField, cmd = m.InputField.Update(msg)
	}

	return m, cmd
}

func performAction(choice int) tea.Cmd {
	return func() tea.Msg {
		switch choice {
		case 0:
			return createSSHKey("id_rsa", "")
		case 1:
			setGitConfig()
		case 2:
			cloneRepo()
		case 3:
			commitAndSync()
		case 4:
			return tea.Quit
		}
		return nil
	}
}

// Ici tu penses à ton interface totale
// Qui s'adapte au model
// Oublie le string builder
func (m model) View() string {
	var output string

	output += titleASCII + "\n"

	if m.ShowInputField {
		output += fmt.Sprintf(
			"Input ?\n\n%s\n\n%s",
			inputStyle.Render(m.InputField.View()),
			subtleStyle.Render("(esc to quit)"),
		) + "\n"
	}

	if m.InputField.Value() != "" && !m.ShowInputField {
		output += fmt.Sprintf("Value : %s", m.InputField.Value())
	}
	output += fmt.Sprintf("\n\n%s", choicesView(m))

	return mainStyle.Render(output)
}

// The first view, where you're choosing a task
func choicesView(m model) string {
	c := m.Choice

	tpl := "%s\n\n"
	tpl += subtleStyle.Render(" j/k, up/down: select ") +
		subtleStyle.Render(" enter: choose") +
		subtleStyle.Render(" q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		renderOption(keywordStyle.Render(optionCreateSSHKey), c == 0),
		renderOption(keywordStyle.Render(optionSetGlobalGitConfig), c == 1),
		renderOption(keywordStyle.Render(optionCloneGitHubRepo), c == 2),
	)

	return fmt.Sprintf(tpl, choices)
}

func renderOption(label string, selected bool) string {
	if selected {
		return selectedOptionStyle.Render(bracketStyle.Render("[x] ") + checkboxStyle.Render(label))
	}
	return optionStyle.Render(bracketStyle.Render("[ ] ") + label)
}

func runCommand(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}

func createSSHKey(filePath, passphrase string) tea.Cmd {
	return func() tea.Msg {
		args := []string{"-t", "rsa", "-b", "4096", "-C", "your_email@example.com"}
		if passphrase != "" {
			args = append(args, "-N", passphrase)
		}
		args = append(args, "-f", filePath)

		cmd := exec.Command("ssh-keygen", args...)
		cmd.Stdin = os.Stdin
		runCommand(cmd)

		fmt.Println(ticksStyle.Render("SSH key created successfully!"))
		return nil
	}
}

func setGitConfig() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter global username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter global email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	runCommand(exec.Command("git", "config", "--global", "user.name", username))
	runCommand(exec.Command("git", "config", "--global", "user.email", email))
}

func cloneRepo() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter GitHub repository URL: ")
	repoURL, _ := reader.ReadString('\n')
	repoURL = strings.TrimSpace(repoURL)

	runCommand(exec.Command("git", "clone", repoURL))
}

func commitAndSync() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter commit message: ")
	commitMessage, _ := reader.ReadString('\n')
	commitMessage = strings.TrimSpace(commitMessage)

	runCommand(exec.Command("git", "add", "."))
	runCommand(exec.Command("git", "commit", "-m", commitMessage))
	runCommand(exec.Command("git", "push", "origin", "main"))
}

func main() {
	// Check the operating system
	if runtime.GOOS == "windows" {
		fmt.Println("Note: This script is running on Windows. Some commands might differ.")
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
