package tui

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"fmt"
	"sort"

	"github.com/brittonhayes/therapy"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#ffffff"))

var bannerStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#477be4")).
	Foreground(lipgloss.Color("#ffffff"))

var focusedStyle = lipgloss.NewStyle().Width(120).Padding(1).Faint(true)

type model struct {
	banner    string
	title     string
	statement string
	Viewport  viewport.Model
	Table     table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Sequence(
				tea.ExitAltScreen,
				tea.ClearScreen,
				tea.Quit,
			)
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m model) bannerView() string {
	return bannerStyle.Render(m.banner)
}

func (m model) footerView() string {
	selection := fmt.Sprintf("%s - %s\n%s", m.Table.SelectedRow()[0], m.Table.SelectedRow()[1], m.Table.SelectedRow()[2])
	return focusedStyle.Render(selection)
}

func (m model) bodyView() string {
	return baseStyle.Render(m.Table.View())
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n", m.bannerView(), m.bodyView(), m.footerView())
}

func Run(therapists []therapy.Therapist) error {

	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Phone", Width: 15},
		{Title: "Credentials", Width: 40},
	}
	rows := []table.Row{}

	sort.SliceStable(therapists, func(i, j int) bool {
		return therapists[i].Title < therapists[j].Title
	})

	for _, t := range therapists {
		if len(t.Phone) == 0 {
			t.Phone = "N/A"
		}

		rows = append(rows, table.Row{
			t.Title,
			t.Phone,
			t.Credentials,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#f5f7f9")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#477be4")).
		Italic(true).
		Bold(true)

	t.SetStyles(s)

	m := model{Table: t}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}
