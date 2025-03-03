package model

type Model struct {
	Choice        int
	CurrentOption int
	Options       []string
	SubMenu       []string
	MenuStack     [][]string
	InputMode     bool
	InputField    string
	InputValue    string
	FilePath      string
	Passphrase    string
	Step          int
}

const (
	OptionCreateSSHKey       = "Create SSH Key"
	OptionSetGlobalGitConfig = "Set Global Git Config"
	OptionCloneGitHubRepo    = "Clone GitHub Repository"
	OptionCommitAndSync      = "Commit and Sync Changes"
	OptionExit               = "Exit"
	OptionBack               = "Back"
	OptionCancel             = "Cancel"
	OptionFilePath           = "Enter file path"
	OptionPassphrase         = "Enter passphrase"
	OptionGenerateKeys       = "Generate Keys"
)

func InitialModel() Model {
	return Model{
		Options:   []string{OptionCreateSSHKey, OptionSetGlobalGitConfig, OptionCloneGitHubRepo, OptionCommitAndSync, OptionExit},
		MenuStack: [][]string{},
		Step:      0,
	}
}
