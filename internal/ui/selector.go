package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	promptStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

type model struct {
	items    []string
	filtered []string
	cursor   int
	input    string
	selected string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor]
			}
			m.quitting = true
			return m, tea.Quit
		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.filter()
			}
		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
				m.filter()
			}
		}
	}
	return m, nil
}

func (m *model) filter() {
	if m.input == "" {
		m.filtered = m.items
		m.cursor = 0
		return
	}

	matches := fuzzy.Find(m.input, m.items)
	m.filtered = make([]string, len(matches))
	for i, match := range matches {
		m.filtered[i] = match.Str
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = 0
	}
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Prompt
	b.WriteString(promptStyle.Render("> "))
	b.WriteString(m.input)
	b.WriteString("\n\n")

	// Items
	maxItems := 10
	start := 0
	if m.cursor >= maxItems {
		start = m.cursor - maxItems + 1
	}

	for i := start; i < len(m.filtered) && i < start+maxItems; i++ {
		item := m.filtered[i]
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("→ " + item))
		} else {
			b.WriteString(normalStyle.Render("  " + item))
		}
		b.WriteString("\n")
	}

	if len(m.filtered) == 0 {
		b.WriteString(normalStyle.Render("  No matches"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(normalStyle.Render("↑/↓: navigate • enter: select • esc: cancel"))

	return b.String()
}

// FuzzySelect presents an interactive fuzzy finder and returns the selected item
func FuzzySelect(items []string) (string, error) {
	m := model{
		items:    items,
		filtered: items,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run selector: %w", err)
	}

	result := finalModel.(model)
	return result.selected, nil
}
