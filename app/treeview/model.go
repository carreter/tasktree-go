package treeview

import (
	"fmt"
	"github.com/carreter/tasktree-go/pkg/tasktree"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type Model struct {
	taskTree *tasktree.TaskTree
	root     *uuid.UUID
}

func NewModel(taskTree *tasktree.TaskTree) Model {
	return Model{
		taskTree: taskTree,
		root:     nil,
	}
}

func (m Model) Root(id uuid.UUID) {
	m.root = &id
}

func (m Model) ClearRoot() {
	m.root = nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	tree, err := RenderTaskTree(m.taskTree, 4)
	if err != nil {
		return fmt.Sprintf("could not render tree: %v", err)
	}

	return tree
}
