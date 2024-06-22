package treeview

import (
	"fmt"
	"github.com/carreter/tasktree-go/task"
	"github.com/carreter/tasktree-go/tasktree"
	"github.com/charmbracelet/lipgloss/tree"
	"strings"
)

func RenderTaskTreeFromRoot(taskTree *tasktree.TaskTree, rootId task.Id, depth int) (string, error) {
	root, exists := taskTree.GetTask(rootId)
	if !exists {
		return "", fmt.Errorf("task %v not found", rootId)
	}

	t := tree.New().Root(root.Name)

	if depth == 0 {
		return t.String(), nil
	}

	subtasks, err := taskTree.GetDirectSubtasksOf(rootId)
	if err != nil {
		return "", nil
	}
	for _, subtask := range subtasks {
		subTree, err := RenderTaskTreeFromRoot(taskTree, subtask.Id, depth-1)
		if err != nil {
			return "", err
		}
		t.Child(subTree)
	}

	return t.String(), nil
}

func RenderTaskTree(taskTree *tasktree.TaskTree, depth int) (string, error) {
	trees := make([]string, len(taskTree.GetRootTasks()))
	for i, rootTask := range taskTree.GetRootTasks() {
		newTree, err := RenderTaskTreeFromRoot(taskTree, rootTask.Id, depth-1)
		if err != nil {
			return "", err
		}
		trees[i] = newTree
	}

	return strings.Join(trees, "\n"), nil
}
