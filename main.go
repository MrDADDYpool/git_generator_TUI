package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
)

var (
	titleASCII = `
   ####    ####    ######   ##   ##  ##   ##  ######            ##   ##    ##     ##   ##    ##       ####   #######  ######
  ##  ##    ##     # ## #   ##   ##  ##   ##   ##  ##           ### ###   ####    ###  ##   ####     ##  ##   ##   #   ##  ##
 ##         ##       ##     ##   ##  ##   ##   ##  ##           #######  ##  ##   #### ##  ##  ##   ##        ## #     ##  ##
 ##         ##       ##     #######  ##   ##   #####            #######  ##  ##   ## ####  ##  ##   ##        ####     #####
 ##  ###    ##       ##     ##   ##  ##   ##   ##  ##           ## # ##  ######   ##  ###  ######   ##  ###   ## #     ## ##
  ##  ##    ##       ##     ##   ##  ##   ##   ##  ##           ##   ##  ##  ##   ##   ##  ##  ##    ##  ##   ##   #   ##  ##
   #####   ####     ####    ##   ##   #####   ######            ##   ##  ##  ##   ##   ##  ##  ##     #####  #######  #### ##
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
				}
			} else if m.options[m.currentOption] == optionCancel {
				return m, tea.Quit
			} else if m.options[m.currentOption] == optionCreateSSHKey {
				m.menuStack = append(m.menuStack, m.options)
				m.options = []string{optionFilePath, optionPassphrase, optionBack, optionCancel}
				m.currentOption = 0
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
				// Process the input value
				if m.inputField == optionFilePath {
					// Handle file path input
				} else if m.inputField == optionPassphrase {
					// Handle passphrase input
				}
			} else {
				m.inputMode = true
				m.inputField = m.options[m.currentOption]
				m.inputValue = ""
			}
		case "esc":
			m.inputMode = false
			m.inputValue = ""
		}
	}
	return m, nil
}

func performAction(choice int) tea.Cmd {
	return func() tea.Msg {
		switch choice {
		case 1:
			createSSHKey()
		case 2:
			setGitConfig()
		case 3:
			cloneRepo()
		case 4:
			commitAndSync()
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

func createSSHKey() {
	reader := bufio.NewReader(os.Stdin)
	var filePath, passphrase string

	// Demander le chemin de sauvegarde
	fmt.Print("Enter file path to save the SSH key (e.g., /home/user/.ssh/id_rsa): ")
	filePath, _ = reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	// Demander la phrase secrète
	fmt.Print("Enter a passphrase for the SSH key (leave empty for no passphrase): ")
	passphrase, _ = reader.ReadString('\n')
	passphrase = strings.TrimSpace(passphrase)

	args := []string{"-t", "rsa", "-b", "4096", "-C", "your_email@example.com"}
	if passphrase != "" {
		args = append(args, "-N", passphrase)
	}
	args = append(args, "-f", filePath)

	cmd := exec.Command("ssh-keygen", args...)
	cmd.Stdin = os.Stdin
	runCommand(cmd)

	fmt.Println("SSH key created successfully!")
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
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
