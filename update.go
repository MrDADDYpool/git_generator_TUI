package update

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	model.go
	commands.go
)

func Update(msg tea.Msg, m model.Model) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.CurrentOption > 0 {
				m.CurrentOption--
			}
		case "down", "j":
			if m.CurrentOption < len(m.Options)-1 {
				m.CurrentOption++
			}
		case " ":
			if m.Options[m.CurrentOption] == model.OptionBack {
				if len(m.MenuStack) > 0 {
					m.MenuStack = m.MenuStack[:len(m.MenuStack)-1]
					if len(m.MenuStack) == 0 {
						m.Options = []string{model.OptionCreateSSHKey, model.OptionSetGlobalGitConfig, model.OptionCloneGitHubRepo, model.OptionCommitAndSync, model.OptionExit}
					} else {
						m.Options = m.MenuStack[len(m.MenuStack)-1]
					}
					m.InputMode = false
					m.Step = 0
				}
			} else if m.Options[m.CurrentOption] == model.OptionCancel {
				return m, tea.Quit
			} else if m.Options[m.CurrentOption] == model.OptionCreateSSHKey {
				m.MenuStack = append(m.MenuStack, m.Options)
				m.Options = []string{model.OptionFilePath, model.OptionPassphrase, model.OptionGenerateKeys, model.OptionBack, model.OptionCancel}
				m.CurrentOption = 0
				m.Step = 1
			} else if m.Options[m.CurrentOption] == model.OptionGenerateKeys {
				return m, commands.CreateSSHKey(m.FilePath, m.Passphrase)
			} else {
				m.MenuStack = append(m.MenuStack, m.Options)
				m.Options = m.SubMenu
				m.CurrentOption = 0
				return m, commands.PerformAction(m.Choice)
			}
		case "enter":
			if m.InputMode {
				m.InputMode = false
				m.InputValue = strings.TrimSpace(m.InputValue)
				if m.Step == 1 {
					m.FilePath = m.InputValue
					m.Step = 2
					m.InputField = model.OptionPassphrase
					m.InputValue = ""
				} else if m.Step == 2 {
					m.Passphrase = m.InputValue
					m.Step = 3
					return m, commands.CreateSSHKey(m.FilePath, m.Passphrase)
				}
			} else {
				m.InputMode = true
				m.InputField = m.Options[m.CurrentOption]
				m.InputValue = ""
			}
		case "esc":
			m.InputMode = false
			m.InputValue = ""
		case "ctrl+u":
			if m.InputMode {
				m.InputValue = ""
			}
		case "ctrl+w":
			if m.InputMode && len(m.InputValue) > 0 {
				fields := strings.Fields(m.InputValue)
				if len(fields) > 0 {
					m.InputValue = strings.Join(fields[:len(fields)-1], " ")
				}
			}
		default:
			if m.InputMode {
				m.InputValue += msg.String()
			}
		}
	}
	return m, nil
}