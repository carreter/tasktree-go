package command

import (
	"fmt"
	"github.com/carreter/tasktree-go/app"
	"github.com/carreter/tasktree-go/pkg/task"
	"github.com/carreter/tasktree-go/pkg/tasktree"
	"github.com/carreter/tasktree-go/pkg/util"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"unicode"
)

type Command interface {
	Call(...string) (string, error)
	Usage() string
	Name() string
}

type Model struct {
	ctx       *app.Context
	taskTree  *tasktree.TaskTree
	focused   bool
	textInput textinput.Model
	errorMsg  string
	commands  []Command
}

func New(taskTree *tasktree.TaskTree) Model {
	textInput := textinput.New()
	return Model{
		taskTree:  taskTree,
		textInput: textInput,
	}
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

	switch args[0] {
	case "q", "quit":
		return tea.Quit
	case "add":
		if len(args) < 2 {
			m.errorMsg = "usage: add <task name> [<task description]"
		}

		newTask := task.Task{Id: task.Id(args[1]), Name: args[1]}
		if len(args) > 2 {
			newTask.Description = strings.Trim(strings.Join(args[2:], " "), "\"")
		}

		err := m.taskTree.AddTask(newTask)
		if err != nil {
			m.errorMsg = fmt.Sprintf("failed to add task: %v", err)
		}
	case "add-subtask":
		if len(args) < 3 {
			m.errorMsg = "usage: add-subtask <parent task id> <subtask name> [<subtask description]"
			return nil
		}

		newTask := task.Task{Id: task.Id(args[2]), Name: args[2]}
		if len(args) > 3 {
			newTask.Description = strings.Trim(strings.Join(args[3:], " "), "\"")
		}

		err := m.taskTree.AddTask(newTask)
		if err != nil {
			m.errorMsg = fmt.Sprintf("failed to add task: %v", err)
			return nil
		}
		err = m.taskTree.MarkSubtask(task.Id(args[1]), newTask.Id)
		if err != nil {
			m.errorMsg = fmt.Sprintf("failed to add subtask: %v", err)
			return nil
		}
	case "delete":
		if len(args) != 2 {
			m.errorMsg = "usage: delete <task id>"
			return nil
		}
		err := m.taskTree.DeleteTask(task.Id(args[1]))
		if err != nil {
			m.errorMsg = fmt.Sprintf("failed to delete task: %v", err)
			return nil
		}
	default:
		m.errorMsg = fmt.Sprintf("unknown command: %s", args[0])
	}

	return nil
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
