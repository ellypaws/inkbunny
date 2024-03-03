package gui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ellypaws/inkbunny/utils"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type listModel struct {
	list   list.Model
	active bool
}

func (m listModel) Focus() tea.Model {
	m.active = true
	return m
}

func (m listModel) Blur() tea.Model {
	m.active = false
	return m
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Index() int {
	return m.list.Index()
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.active {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.active && msg.String() == "enter" {
			return m, utils.Wrap(m.list.SelectedItem().(item))
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if !m.active {
		return ""
	}
	return docStyle.Render(m.list.View())
}

func initialMenu(items []list.Item) listModel {
	m := listModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Menu"

	return m
}
