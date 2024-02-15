package gui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"encoding/json"
	"fmt"
	"inkbunny/entities"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type loginForm struct {
	user       *entities.Login
	focusIndex int
	submitted  bool
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func InitialModel(user *entities.Login) loginForm {
	m := loginForm{
		user:   user,
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.CharLimit = 64
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m loginForm) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.submitted {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				for i := range m.inputs {
					m.inputs[i].Blur()
				}

				m.submitted = true

				m.user.Username = m.inputs[0].Value()
				m.user.Password = m.inputs[1].Value()

				login := m.login(m.user)

				return m, login
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

// Wrap casts a message into a tea.Cmd
func Wrap(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

// login sends a login request to the server using an *entities.Login
// It mutates the loginForm and returns *entities.Login
func (m *loginForm) login(user *entities.Login) tea.Cmd {
	if user.Username == "" {
		user.Username = "guest"
	} else if user.Password == "" {
		return Wrap(fmt.Errorf("username is set but password is empty"))
	}
	resp, err := http.PostForm(entities.BaseURL+"login.php", url.Values{"username": {user.Username}, "password": {user.Password}})
	if err != nil {
		return Wrap(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Wrap(err)
	}

	if err = json.Unmarshal(body, user); err != nil {
		return Wrap(err)
	}

	return Wrap(user)
}

func (m *loginForm) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		switch i {
		case 0:
			m.user.Username = m.inputs[i].Value()
		case 1:
			m.user.Password = m.inputs[i].Value()
		}
	}

	return tea.Batch(cmds...)
}

func (m loginForm) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s", *button)

	return b.String()
}
