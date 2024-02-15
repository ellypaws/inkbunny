package main

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"inkbunny/entities"
	"inkbunny/gui"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

const (
	baseURL = "https://inkbunny.net/api_"
)

// Function to get watchlist for a given user
func getWatchlist(sid string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%swatchlist.php?sid=%s&user_id=%s", baseURL, sid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var watchResp entities.WatchlistResponse
	if err := json.Unmarshal(body, &watchResp); err != nil {
		return nil, err
	}

	var usernames []string
	for _, watch := range watchResp.Watches {
		usernames = append(usernames, watch.Username)
	}

	return usernames, nil
}

// Function to find mutual elements in two slices
func findMutual(a, b []string) []string {
	var mutual []string
	for _, val := range a {
		if slices.Contains(b, val) {
			mutual = append(mutual, val)
		}
	}
	return mutual
}

func changeRating(sid string) error {
	resp, err := http.PostForm(baseURL+"userrating.php", url.Values{
		"sid":    {sid},
		"tag[2]": {"yes"},
		"tag[3]": {"yes"},
		"tag[4]": {"yes"},
		"tag[5]": {"yes"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResp entities.Login
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return err
	}

	if loginResp.Sid != sid {
		return fmt.Errorf("session ID changed after rating change, expected: [%s], got: [%s]", sid, loginResp.Sid)
	}

	return nil
}

func logout(sid string) error {
	resp, err := http.PostForm(baseURL+"logout.php", url.Values{"sid": {sid}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var logoutResp entities.LogoutResponse
	if err = json.Unmarshal(body, &logoutResp); err != nil {
		return err
	}

	if logoutResp.Logout != "success" {
		return fmt.Errorf("logout failed, response: %s", logoutResp.Logout)
	}

	return nil
}

func getUserID(username string) (entities.User, error) {
	resp, err := http.Get(fmt.Sprintf("%susername_autosuggest.php?username=%s", baseURL, username))
	if err != nil {
		return entities.User{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entities.User{}, err
	}

	var user entities.User
	if err := json.Unmarshal(body, &user); err != nil {
		return entities.User{}, err
	}

	return user, nil
}

type model struct {
	user entities.Login
	l    tea.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		log.Println("Error:", msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		m.l, cmd = m.l.Update(msg)
		return m, cmd
	case *entities.Login:
		if msg == nil {
			log.Println("Login message is nil")
			return m, tea.Quit
		}
		m.user = *msg
		m, cmd = m.login()
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	view := strings.Builder{}
	view.WriteString("Inkbunny CLI\n\n")
	view.WriteString(m.l.View())
	switch {
	case m.user.Sid != "":
		view.WriteString(fmt.Sprintf("\n\nLogged in as [%s] with session ID [%s]", m.user.Username, m.user.Sid))

		watchlist, err := getWatchlist(m.user.Sid)
		if err != nil {
			log.Println("Error getting watchlist:", err)
		}

		view.WriteString(fmt.Sprintf("\nWatch list: %v", watchlist))
	}
	return view.String()
}

func initialModel() model {
	return model{
		l: gui.InitialModel(&entities.Login{}),
	}
}

func (m model) login() (model, tea.Cmd) {
	user := &m.user
	if user.Username == "" {
		user.Username = "guest"
	} else if user.Password == "" {
		fmt.Errorf("username is set but password is empty")
	}
	resp, err := http.PostForm(baseURL+"login.php", url.Values{"username": {user.Username}, "password": {user.Password}})
	if err != nil {
		return m, wrap(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return m, wrap(err)
	}

	if err = json.Unmarshal(body, user); err != nil {
		return m, wrap(err)
	}

	return m, nil
}

func wrap(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
