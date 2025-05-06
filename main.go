package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles graphiques
var (
	successStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	borderStyle         = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).BorderForeground(lipgloss.Color("63"))
	titleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	optionStyle         = lipgloss.NewStyle().Padding(1, 2).Foreground(lipgloss.Color("69"))
	selectedOptionStyle = lipgloss.NewStyle().Padding(1, 2).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("63")).Bold(true)
	stepTitleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
)

// ASCII Art du titre
var titleASCII = titleStyle.Render(`
  ________._________________ ___  ____ _____________     _____      _____    _______      _____    _____________________________
 /  _____/|   \__    ___/   |   \|    |   \______   \   /     \    /  _  \   \      \    /  _  \  /  _____/\_   _____/\______   \
/   \  ___|   | |    | /    ~    \    |   /|    |  _/  /  \ /  \  /  /_\  \  /   |   \  /  /_\  \/   \  ___ |    ___  |       _/
/\    \_\  \   | |    | \    Y    /    |  / |    |   \ /    Y    \/    |    \/    |    \/    |    \    \_\  \|        \ |    |   \
 \______  /___| |____|  \___|_  /|______/  |______  / \____|__  /\____|__  /\____|__  /\____|__  /\______  /_______  / |____|_  /
        \/                    \/                  \/          \/         \/         \/         \/        \/        \/         \/
`)

func main() {
	checkDependencies()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erreur : %v\n", err)
		os.Exit(1)
	}
}

// États du menu
type step int

const (
	stepChoice step = iota
	stepGitHub
	stepGitea
	stepCloneRepo
)

// Modèle pour l'interface
type model struct {
	choice   int
	input    textinput.Model
	step     step
	progress progress.Model
}

// FinishedMsg is a custom message to indicate progress completion
type FinishedMsg struct{}

// Implement progress.FinishedMsg to avoid undefined reference
func (FinishedMsg) String() string {
	return "Progress finished"
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "https://gitea.example.com"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40

	// Initialisation de la barre de progression
	progressBar := progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
		progress.WithScaledGradient("green", "blue"),
	)

	return model{
		choice:   0,
		input:    ti,
		step:     stepChoice,
		progress: progressBar,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, startProgress())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.step == stepChoice && m.choice > 0 {
				m.choice--
			}
		case "down":
				if m.step == stepChoice && m.choice < 2 {
					m.choice++
					return m, tea.Batch(startProgress(), func() tea.Msg {
						configureGitService("github.com")
						return FinishedMsg{}
					})
				}
			}
			// Removed unreachable code
			if m.choice == 1 {
					m.step = stepGitea
					m.input.Focus()
				} else {
					m.step = stepCloneRepo
					m.input.Focus()
				}
			}
			// Correctly close the previous if block
			if m.step == stepGitea {
				configureGitService(strings.TrimSpace(m.input.Value()))
				m.input.Reset()
				m.step = stepChoice
			} else if m.step == stepCloneRepo {
				cloneRepository(strings.TrimSpace(m.input.Value()))
				m.input.Reset()
				m.step = stepChoice
			}
		}

	case progress.FrameMsg:
		var progressCmd tea.Cmd
		updatedProgress, progressCmd := m.progress.Update(msg)
		m.progress = updatedProgress
		return m, progressCmd
	}

	if m.step == stepGitea || m.step == stepCloneRepo {
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var output string
	output += titleASCII + "\n"

	switch m.step {
	case stepChoice:
		output += stepTitleStyle.Render("Sélectionnez une option :") + "\n"
		options := []string{"Configurer GitHub", "Configurer Gitea", "Cloner un dépôt"}
		for i, option := range options {
			if i == m.choice {
				output += selectedOptionStyle.Render(option) + "\n"
			} else {
				output += optionStyle.Render(option) + "\n"
			}
		}
		output += "\n↑↓ pour naviguer, Entrée pour sélectionner, q pour quitter."

	case stepGitHub, stepGitea, stepCloneRepo:
		output += stepTitleStyle.Render("Progression :") + "\n"
		output += m.progress.View() + "\n"
		if m.step == stepGitea {
			output += stepTitleStyle.Render("Entrez l'URL de l'instance Gitea :") + "\n"
			output += m.input.View() + "\n(Entrée pour valider)"
		} else if m.step == stepCloneRepo {
			output += stepTitleStyle.Render("Entrez l'URL du dépôt à cloner (HTTPS ou SSH) :") + "\n"
			output += m.input.View() + "\n(Entrée pour valider)"
		}
	}

	return output
}

// Vérifie si les dépendances essentielles sont installées
func checkDependencies() {
	commands := []string{"git", "ssh-keygen"}
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd); err != nil {
			fmt.Printf("Erreur : %s non trouvé. Veuillez l'installer.\n", cmd)
			os.Exit(1)
		}
	}
}

// Configure le service Git
func configureGitService(service string) {
	steps := []string{"Génération de la clé SSH", "Ajout de la clé SSH à l'agent", "Test de la connexion SSH"}

	for i, step := range steps {
		fmt.Printf("%s...\n", step)
		switch i {
		case 0:
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(errorStyle.Render("Erreur : Impossible de récupérer le répertoire utilisateur."))
				return
			}
			sshKeyPath := fmt.Sprintf("%s/.ssh/id_rsa", homeDir)
			generateSSHKey(sshKeyPath)
		case 1:
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(errorStyle.Render("Erreur : Impossible de récupérer le répertoire utilisateur."))
				return
			}
			sshKeyPath := fmt.Sprintf("%s/.ssh/id_rsa", homeDir)
			runCommand(exec.Command("ssh-add", sshKeyPath))
		case 2:
			testSSHConnection(service)
		}
	}
	fmt.Println(successStyle.Render("[100%] Configuration terminée avec succès !"))
}

// Génère une clé SSH
func generateSSHKey(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println(successStyle.Render("Clé SSH déjà existante, aucune action nécessaire."))
		return
	}
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", "your_email@example.com", "-f", filePath, "-N", "")
	runCommand(cmd)
	fmt.Println(successStyle.Render("Clé SSH générée avec succès !"))
}

// Teste la connexion SSH
func testSSHConnection(service string) {
	cmd := exec.Command("ssh", "-T", fmt.Sprintf("git@%s", service))
	runCommand(cmd)
}

// Clone un dépôt Git via HTTPS ou SSH
func cloneRepository(repoURL string) {
	fmt.Println("Clonage du dépôt en cours...")
	cmd := exec.Command("git", "clone", repoURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Erreur lors du clonage du dépôt : %v %s", err, string(output))
		return
	}
	fmt.Println(successStyle.Render("Dépôt cloné avec succès !"))
}

// Exécute une commande système
func runCommand(cmd *exec.Cmd) {
	output, err := cmd.CombinedOutput()
	formattedOutput := borderStyle.Render(string(output))

	if err != nil {
		fmt.Println(errorStyle.Render("Erreur :"), err)
		fmt.Println(formattedOutput)
		return
	}
	fmt.Println(successStyle.Render("Commande exécutée avec succès :"))
	fmt.Println(formattedOutput)
}

func startProgress() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progress.FrameMsg{}
	})
}
