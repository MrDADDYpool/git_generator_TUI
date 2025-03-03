package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choice        int
	currentOption int
	options       []string
	subMenu       []string
	menuStack     [][]string
	inputMode     bool
	inputField    string
	inputValue    string
	filePath      string
	passphrase    string
	step          int
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

func initialModel() model {
	return model{
		options:   []string{optionCreateSSHKey, optionSetGlobalGitConfig, optionCloneGitHubRepo, optionCommitAndSync, optionExit},
		menuStack: [][]string{},
		step:      0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.currentOption > 0 {
				m.currentOption--
			}
		case "down", "j":
			if m.currentOption < len(m.options)-1 {
				m.currentOption++
			}
		case " ":
			if m.options[m.currentOption] == optionBack {
				if len(m.menuStack) > 0 {
					m.menuStack = m.menuStack[:len(m.menuStack)-1]
					if len(m.menuStack) == 0 {
						m.options = []string{optionCreateSSHKey, optionSetGlobalGitConfig, optionCloneGitHubRepo, optionCommitAndSync, optionExit}
					} else {
						m.options = m.menuStack[len(m.menuStack)-1]
					}
					m.inputMode = false
					m.step = 0
				}
			} else if m.options[m.currentOption] == optionCancel {
				return m, tea.Quit
			} else if m.options[m.currentOption] == optionCreateSSHKey {
				m.menuStack = append(m.menuStack, m.options)
				m.options = []string{optionFilePath, optionPassphrase, optionGenerateKeys, optionBack, optionCancel}
				m.currentOption = 0
				m.step = 1
			} else if m.options[m.currentOption] == optionGenerateKeys {
				return m, createSSHKey(m.filePath, m.passphrase)
			} else {
				m.menuStack = append(m.menuStack, m.options)
				m.options = m.subMenu
				m.currentOption = 0
				return m, performAction(m.choice)
			}
		case "enter":
			if m.inputMode {
				m.inputMode = false
				m.inputValue = strings.TrimSpace(m.inputValue)
				if m.step == 1 {
					m.filePath = m.inputValue
					m.step = 2
					m.inputField = optionPassphrase
					m.inputValue = ""
				} else if m.step == 2 {
					m.passphrase = m.inputValue
					m.step = 3
					return m, createSSHKey(m.filePath, m.passphrase)
				}
			} else {
				m.inputMode = true
				m.inputField = m.options[m.currentOption]
				m.inputValue = ""
			}
		case "esc":
			m.inputMode = false
			m.inputValue = ""
		case "ctrl+u":
			if m.inputMode {
				m.inputValue = ""
			}
		case "ctrl+w":
			if m.inputMode && len(m.inputValue) > 0 {
				// Delete the last word
				fields := strings.Fields(m.inputValue)
				if len(fields) > 0 {
					m.inputValue = strings.Join(fields[:len(fields)-1], " ")
				}
			}
		default:
			if m.inputMode {
				m.inputValue += msg.String()
			}
		}
	}
	return m, nil
}

func performAction(choice int) tea.Cmd {
	return func() tea.Msg {
		switch choice {
		case 1:
			// This case is handled in the Update method
		case 2:
			return setGitConfigCmd()
		case 3:
			return cloneRepoCmd()
		case 4:
			return commitAndSyncCmd()
		case 5:
			return tea.Quit
		}
		return nil
	}
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleASCII + "\n\n")

	if m.inputMode {
		b.WriteString(inputStyle.Render(fmt.Sprintf("%s: %s", m.inputField, m.inputValue)) + "\n\n")
	}

	for i, option := range m.options {
		if m.currentOption == i {
			b.WriteString(bracketStyle.Render("[") + cursorStyle.String() + bracketStyle.Render("] ") + selectedOptionStyle.Render(option) + "\n")
		} else {
			b.WriteString(bracketStyle.Render("[ ] ") + optionStyle.Render(option) + "\n")
		}
	}

	return b.String()
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

		fmt.Println("SSH key created successfully!")
		return nil
	}
}

func setGitConfigCmd() tea.Cmd {
	return func() tea.Msg {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter global username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Enter global email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		runCommand(exec.Command("git", "config", "--global", "user.name", username))
		runCommand(exec.Command("git", "config", "--global", "user.email", email))

		return nil
	}
}

func cloneRepoCmd() tea.Cmd {
	return func() tea.Msg {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter GitHub repository URL: ")
		repoURL, _ := reader.ReadString('\n')
		repoURL = strings.TrimSpace(repoURL)

		runCommand(exec.Command("git", "clone", repoURL))
		return nil
	}
}

func commitAndSyncCmd() tea.Cmd {
	return func() tea.Msg {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter commit message: ")
		commitMessage, _ := reader.ReadString('\n')
		commitMessage = strings.TrimSpace(commitMessage)

		runCommand(exec.Command("git", "add", "."))
		runCommand(exec.Command("git", "commit", "-m", commitMessage))
		runCommand(exec.Command("git", "push", "origin", "main"))
		return nil
	}
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
