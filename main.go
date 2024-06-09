package main

import (
	"fmt"
	"os"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor   int
	todos    []string
	selected map[int]struct{}
	adding   bool
	newTodo  string
}

func initialModel() model {
	return model{
		todos:    []string{"Hello, Hamsurang!"},
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
			switch msg.Type {
			case tea.KeyEnter:
				if m.newTodo != "" {
					m.todos = append(m.todos, m.newTodo)
					m.newTodo = ""
					m.adding = false
				}
			case tea.KeyEsc:
				m.adding = false
				m.newTodo = ""
			case tea.KeyBackspace:
				if len(m.newTodo) > 0 {
					_, size := utf8.DecodeLastRuneInString(m.newTodo)
					m.newTodo = m.newTodo[:len(m.newTodo)-size]
				}
			default:

				if msg.Type == tea.KeyRunes {
					m.newTodo += string(msg.Runes)
				}
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
				if m.cursor < len(m.todos)-1 {
					m.cursor++
				}
			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			case "a":
				m.adding = true
				m.newTodo = ""
			case "d":

				if len(m.todos) > 0 {
					m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
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
		return fmt.Sprintf("Add a new task: %s", m.newTodo)
	}

	s := `
Todo List - Use the arrow keys to navigate and [space] to mark tasks as done:
  ┌───────────────────────────────────────────────┐
  │ [a] Add a new task                            │
  │ [d] Delete the current task                   │
  │ [q] Quit                                      │
  └───────────────────────────────────────────────┘

`
	for i, todo := range m.todos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, todo)
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
