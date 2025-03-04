package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/schollz/progressbar"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles graphiques
var (
	titleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	optionStyle         = lipgloss.NewStyle().Padding(1, 2).Foreground(lipgloss.Color("69"))
	selectedOptionStyle = lipgloss.NewStyle().Padding(1, 2).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("63")).Bold(true)
	stepTitleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
)

// ASCII Art du titre
var titleASCII = titleStyle.Render(`
  ________._________________ ___  ____ _____________     _____      _____    _______      _____    _____________________________
 /  _____/|   \__    ___/   |   \|    |   \______   \   /     \    /  _  \   \      \    /  _  \  /  _____/\_   _____/\______   \
/   \  ___|   | |    | /    ~    \    |   /|    |  _/  /  \ /  \  /  /_\  \  /   |   \  /  /_\  \/   \  ___ |    __)_  |       _/
\    \_\  \   | |    | \    Y    /    |  / |    |   \ /    Y    \/    |    \/    |    \/    |    \    \_\  \|        \ |    |   \
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
	choice int
	input  textinput.Model
	step   step
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "https://gitea.example.com"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	return model{choice: 0, input: ti, step: stepChoice}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
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
			}
		case "enter":
			if m.step == stepChoice {
				if m.choice == 0 {
					m.step = stepGitHub
					configureGitService("github.com")
					m.step = stepChoice
				} else if m.choice == 1 {
					m.step = stepGitea
					m.input.Focus()
				} else {
					m.step = stepCloneRepo
					m.input.Focus()
				}
			} else if m.step == stepGitea {
				configureGitService(strings.TrimSpace(m.input.Value()))
				m.input.Reset()
				m.step = stepChoice
			} else if m.step == stepCloneRepo {
				cloneRepository(strings.TrimSpace(m.input.Value()))
				m.input.Reset()
				m.step = stepChoice
			}
		}
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

	case stepGitea:
		output += stepTitleStyle.Render("Entrez l'URL de l'instance Gitea :") + "\n"
		output += m.input.View() + "\n(Entrée pour valider)"

	case stepCloneRepo:
		output += stepTitleStyle.Render("Entrez l'URL du dépôt à cloner (HTTPS ou SSH) :") + "\n"
		output += m.input.View() + "\n(Entrée pour valider)"
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
	bar := progressbar.New(3) // Initialisation de la barre de progression
	steps := []string{"Génération de la clé SSH", "Ajout de la clé SSH à l'agent", "Test de la connexion SSH"}

	for i, step := range steps {
		fmt.Printf("[%d%%] %s...\n", (i+1)*33, step)
		switch i {
		case 0:
			sshKeyPath := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
			generateSSHKey(sshKeyPath)
		case 1:
			sshKeyPath := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
			runCommand(exec.Command("ssh-add", sshKeyPath))
		case 2:
			testSSHConnection(service)
		}
		bar.Add(1) // Mise à jour de la barre de progression
	}
	fmt.Println("[100%] Configuration terminée avec succès !")
}

// Génère une clé SSH
func generateSSHKey(filePath string) {
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", "your_email@example.com", "-f", filePath, "-N", "")
	runCommand(cmd)
	fmt.Println("Clé SSH générée avec succès !")
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
	fmt.Println("Dépôt cloné avec succès !")
}

// Exécute une commande système
func runCommand(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Erreur : %v\n", err)
	}
}
