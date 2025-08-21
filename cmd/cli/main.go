package cli

import (
	"fmt"
	"os"

	"github.com/raghavendrarq/container-dsh/internal/cli"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	fmt.Println("CLI VERSION")
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Renu Madam: ", err)
		os.Exit(1)
	}

}

func initialModel() cli.Model {
	return cli.Model{
		Choices:  []string{"Show containers", "show images", "show stats"},
		Selected: make(map[int]struct{}),
	}
}
