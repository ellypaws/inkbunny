package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"inkbunny/api"
	"inkbunny/gui"
	"log"
	"strings"
)

type model struct {
	user api.Credentials
	l    tea.Model
	p    tea.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		log.Println("Error:", msg)
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		m.l, cmd = m.l.Update(msg)
		return m, cmd
	case *api.Credentials:
		if msg == nil {
			log.Println("Credentials message is nil")
			return m, tea.Quit
		}
		m.user = *msg
		return m, api.Wrap(GetWatchlist{})
	case GetWatchlist:
		watchlist, err := api.GetWatchlist(m.user.Sid)
		if err != nil {
			log.Println("Error getting watchlist:", err)
		}

		m.p = m.p.(gui.Pager).SetContent(strings.Join(watchlist, "\n"))
	}
	m.p, cmd = m.p.Update(msg)
	return m, nil
}

func (m model) View() string {
	view := strings.Builder{}
	view.WriteString("Inkbunny CLI\n\n")
	view.WriteString(m.l.View())
	switch {
	case m.user.Sid != "":
		view.WriteString(fmt.Sprintf("\n\nLogged in as [%s] with session ID [%s]\n", m.user.Username, m.user.Sid))
		view.WriteString(m.p.View())
	}
	return view.String()
}

func initialModel() model {
	return model{
		l: gui.InitialModel(&api.Credentials{}),
		p: gui.NewPager("Fetching watchlist..."),
	}
}

type GetWatchlist struct{}

func main() {
	if _, err := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
