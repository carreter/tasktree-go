package app

import (
	"github.com/carreter/tasktree-go/app/commandview"
	"github.com/carreter/tasktree-go/app/treeview"
	"github.com/carreter/tasktree-go/pkg/tasktree"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focus int

const (
	treeViewFocus focus = iota
	commandFocus
)

type Model struct {
	taskTree *tasktree.TaskTree

	commandView      commandview.Model
	commandViewStyle lipgloss.Style

	treeView      treeview.Model
	treeViewStyle lipgloss.Style

	width  int
	height int

	focus focus
}

func NewModel(taskTree *tasktree.TaskTree) Model {
	commandView := textinput.New()
	commandView.Prompt = ""

	return Model{
		taskTree:    taskTree,
		commandView: commandview.New(taskTree),
		treeView:    treeview.NewModel(taskTree),
		focus:       treeViewFocus,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var focusedCmd tea.Cmd
	switch m.focus {
	case treeViewFocus:
		var newTreeView tea.Model
		newTreeView, focusedCmd = m.treeView.Update(msg)
		m.treeView = newTreeView.(treeview.Model)
	case commandFocus:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "enter":
				m.focus = treeViewFocus
			}
		}

		var newCommandView tea.Model
		newCommandView, focusedCmd = m.commandView.Update(msg)
		m.commandView = newCommandView.(commandview.Model)
	}

	var globalCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			globalCmd = tea.Quit
		case "q":
			if m.focus != commandFocus {
				globalCmd = tea.Quit
			}
		case ":":
			if m.focus != commandFocus {
				m.focus = commandFocus
				m.commandView.Focus()
			}
		}
	}

	return m, tea.Batch(focusedCmd, globalCmd)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		//fmt.Sprintf("%+v", m.taskTree),
		m.treeViewStyle.Render(m.treeView.View()),
		m.commandViewStyle.Render(m.commandView.View()),
	)
}
