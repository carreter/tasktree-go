package command

import (
	"fmt"
	"github.com/carreter/tasktree-go/app"
	"github.com/carreter/tasktree-go/app/views/command/commands"
	"github.com/carreter/tasktree-go/pkg/util"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"unicode"
)

type Command interface {
	Run(*app.Context, ...string) (string, string)
	Usage() string
	Name() string
}

type Model struct {
	ctx *app.Context

	focused   bool
	textInput textinput.Model
	errorMsg  string
	outMsg    string

	commands map[string]Command
}

func New(ctx *app.Context) Model {
	textInput := textinput.New()
	model := Model{
		textInput: textInput,
		ctx:       ctx,
	}
	model.RegisterCommand(commands.DeleteCommand{})
	model.RegisterCommand(commands.AddCommand{})
	model.RegisterCommand(commands.AddSubtaskCommand{})

	return model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmd = m.CallCommand()
			m.textInput.Reset()
			m.Focus()
		case "esc":
			m.textInput.Reset()
			m.Focus()
			return m, nil
		}

	}

	var textInputCmd tea.Cmd
	m.textInput, textInputCmd = m.textInput.Update(msg)

	return m, tea.Batch(textInputCmd, cmd)
}

func (m Model) View() string {
	if !m.focused {
		m.textInput.Prompt = ""
		if m.errorMsg != "" {
			m.textInput.Placeholder = m.errorMsg
			m.textInput.PlaceholderStyle = m.textInput.PlaceholderStyle.Bold(true).Foreground(lipgloss.Color("1"))
		} else if m.outMsg != "" {
			m.textInput.Placeholder = m.outMsg
		} else {
			m.textInput.Placeholder = "Type \":\" to enter command mode"
		}
	} else {
		m.textInput.Prompt = ":"
		m.textInput.Placeholder = ""
	}
	return m.textInput.View()
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() tea.Cmd {
	m.focused = !m.focused
	return m.textInput.Focus()
}

func (m *Model) CallCommand() tea.Cmd {
	args := parseRawArgs(m.textInput.Value())

	if len(args) == 0 {
		m.errorMsg = ""
		return nil
	}

	if args[0] == "q" || args[0] == "quit" {
		return tea.Quit
	}

	command, exists := m.commands[args[0]]
	if !exists {
		m.errorMsg = fmt.Sprintf("unknown command: %s", args[0])
		return nil
	}

	m.outMsg, m.errorMsg = command.Run(m.ctx, args...)

	return nil
}

func (m *Model) RegisterCommand(cmd Command) {

}

func parseRawArgs(raw string) []string {
	quoted := false
	args := strings.FieldsFunc(raw, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}

		return unicode.IsSpace(r) && !quoted
	})
	return util.Map(args, func(s string) string { return strings.Trim(s, "\"") })
}
