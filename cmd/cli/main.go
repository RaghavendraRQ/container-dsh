package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func main() {
	fmt.Println("CLI VERSION")

}

func initialModel() Model {
	return Model{
		choices:  []string{"Show containers", "show images", "show stats"},
		selected: make(map[int]struct{}),
	}
}
