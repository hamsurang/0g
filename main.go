package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor    int
	choices   []string
	selected  map[int]struct{}
	adding    bool
	newChoice string
}

func initialModel() model {
	return model{
		choices:  []string{"Hello, Hamsurang!"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.adding {
			switch msg.String() {
			case "enter":
				if m.newChoice != "" {
					m.choices = append(m.choices, m.newChoice)
					m.newChoice = ""
					m.adding = false
				}
			case "esc":
				m.adding = false
				m.newChoice = ""
			default:
				m.newChoice += msg.String()
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			case "0":
				m.adding = true
				m.newChoice = ""
			case "g":
				if len(m.choices) > 0 {
					m.choices = append(m.choices[:m.cursor], m.choices[m.cursor+1:]...)
					delete(m.selected, m.cursor)
					if m.cursor > 0 {
						m.cursor--
					}
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.adding {
		return fmt.Sprintf("Add a new Task: %s", m.newChoice)
	}

	s := `
Todo List - Use the arrow keys to navigate and [space] to mark tasks as done:
	┌───────────────────────────────────────────────┐
	│ [0] Add a new task                            │
	│ [g] Delete the current task                   │
	│ [q] Quit                                      │
	└───────────────────────────────────────────────┘

`
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
