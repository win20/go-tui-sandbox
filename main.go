package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotui/internal/tui"
	"log"
	"os"
	"strings"
)

const (
	width       = 96
	columnWidth = 30
)

type command struct {
	disabled bool
	name     string
}

type model struct {
	stateDescription string
	stateStatus      tui.StatusBarState
	commands         []command
	cursor           int
	secondListHeader string
	secondListValues []string
}

func initialModel() model {
	return model{
		stateDescription: "Initializing...",
		commands: []command{
			{name: "Set user"},
			{name: "Fetch token", disabled: true},
			{name: "Other"},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				if m.commands[m.cursor - 1].disabled {
					m.cursor--
				}
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.commands)-1 {
				if m.commands[m.cursor + 1].disabled {
					m.cursor++
				}
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	doc := &strings.Builder{}

	tui.RenderStatusBar(doc, tui.NewStatusBarProps(&tui.StatusBarProps{
		Description: m.stateDescription,
		User:        "NONE",
		StatusState: tui.StatusBarStateBlue,
		Width:       width,
	}))


	tui.RenderTitleRow(width, doc, tui.TitleRowProps{Title: "GO TUI example"})
	doc.WriteString("\n\n")

	renderLists(doc, m)

	// Footer
	doc.WriteString("Press q to quit.")
	doc.WriteString("\n")

	// Send UI for rendering
	return doc.String()
}

func renderLists(doc *strings.Builder, m model) {
	var items []tui.Item
	for _, c := range m.commands {
		items = append(items, tui.Item{
			Value:    c.name,
			Disabled: c.disabled,
		})
	}	

	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		tui.RenderListCommands(doc, &tui.ListProps{
			Items:    items,
			Selected: m.cursor,
		}),
		tui.RenderListDisplay(m.secondListHeader, m.secondListValues),
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists))
	doc.WriteString("\n\n")
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")

		if err != nil {
			log.Fatalf("failed setting the debug log file: %v", err)
		}

		defer f.Close()
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("TUI run error: %v", err)
	}
}

