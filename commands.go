package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func runCommand(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}

func CreateSSHKey(filePath, passphrase string) tea.Cmd {
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

func PerformAction(choice int) tea.Cmd {
	return func() tea.Msg {
		switch choice {
		case 2:
			return SetGitConfigCmd()
		case 3:
			return CloneRepoCmd()
		case 4:
			return CommitAndSyncCmd()
		case 5:
			return tea.Quit
		}
		return nil
	}
}

func SetGitConfigCmd() tea.Cmd {
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

func CloneRepoCmd() tea.Cmd {
	return func() tea.Msg {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter GitHub repository URL: ")
		repoURL, _ := reader.ReadString('\n')
		repoURL = strings.TrimSpace(repoURL)

		runCommand(exec.Command("git", "clone", repoURL))
		return nil
	}
}

func CommitAndSyncCmd() tea.Cmd {
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
