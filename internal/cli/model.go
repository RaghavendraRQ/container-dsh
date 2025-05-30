package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Choices  []string
	Cursor   int
	Selected map[int]struct{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "j", "down":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "k", "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			_, ok := m.Selected[m.Cursor]
			if ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = struct{}{} // TODO: Show the Output of the Selected option
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "Select options:\n\n"

	for i, choice := range m.Choices {
		Cursor := " "
		if m.Cursor == i {
			Cursor = ">"
		}

		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", Cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}
