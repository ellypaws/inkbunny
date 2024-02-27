package gui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"inkbunny/api"
	"inkbunny/utils"
	"log"
	"strings"
	"time"
)

type model struct {
	user api.Credentials
	l    tea.Model
	p    tea.Model
	menu tea.Model
}

func (m model) Init() tea.Cmd {
	if m.user.Sid != "" {
		return utils.Wrap(ShowMenu{})
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		log.Println("Error:", msg)
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.user.Sid != "" {
				if err := m.user.Logout(); err != nil {
					return m, utils.Wrap(err)
				}
				log.Println("Logged out")
				time.Sleep(1 * time.Second)
			}
			return InitialModel(""), nil
		}
	case *api.Credentials:
		if msg == nil {
			log.Println("Credentials message is nil")
			return m, tea.Quit
		}
		m.user = *msg
		return m, utils.Wrap(GetWatchlist{})
	case GetWatchlist:
		watchlist, err := api.GetWatchlist(m.user)
		if err != nil {
			log.Println("Error getting watchlist:", err)
		}

		m.p = m.p.(Pager).SetContent(strings.Join(watchlist, "\n"))
		m.chooseFocus(&m.p)
	case ShowMenu:
		m.chooseFocus(&m.menu)
	}
	m, cmds := m.propagate(msg)
	return m, tea.Batch(cmds...)
}

type Focusable interface {
	tea.Model
	Focus() tea.Model
	Blur() tea.Model
}

func (m *model) chooseFocus(model *tea.Model) {
	m.l = m.l.(loginForm).Blur()
	m.p = m.p.(Pager).Blur()
	m.menu = m.menu.(listModel).Blur()

	*model = (*model).(Focusable).Focus()
}

func (m model) propagate(msg tea.Msg) (model, []tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.l, cmd = m.l.Update(msg)
	cmds = append(cmds, cmd)
	m.p, cmd = m.p.Update(msg)
	cmds = append(cmds, cmd)
	m.menu, cmd = m.menu.Update(msg)
	cmds = append(cmds, cmd)
	return m, cmds
}

func (m model) View() string {
	view := strings.Builder{}
	view.WriteString("Inkbunny CLI\n\n")
	view.WriteString(m.l.View())
	switch {
	case m.user.Sid != "":
		view.WriteString(fmt.Sprintf("\n\nLogged in as [%s] with session ID [%s]\n", m.user.Username, m.user.Sid))
		view.WriteString(m.menu.View())
	}
	return view.String()
}

func InitialModel(sid string) model {
	user := api.Credentials{Sid: sid}
	items := []list.Item{
		item{title: "Watchlist", desc: "View your watchlist"},
		item{title: "Logout", desc: "Log out of your account"},
		item{title: "Submissions", desc: "View your submissions"},
		item{title: "Search", desc: "Search for submissions"},
	}
	return model{
		user: user,
		l:    initLoginForm(&user),
		p:    newPager("Fetching watchlist..."),
		menu: initialMenu(items),
	}
}

type GetWatchlist struct{}
type ShowMenu struct{}
