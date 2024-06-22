package tasktree

import (
	"fmt"
	"github.com/carreter/tasktree-go/pkg/task"
	"github.com/carreter/tasktree-go/pkg/util"
	"sync"
)

// A TaskTree organizes the relationships between tasks. Thread-safe.
type TaskTree struct {
	rwMu sync.RWMutex

	tasks map[task.Id]task.Task

	roots []task.Id // tasks that don't have parents

	subtasks  map[task.Id][]task.Id // map from tasks to their subtasks
	subtaskOf map[task.Id]task.Id   // map from subtasks to their parent tasks

	blocks    map[task.Id][]task.Id // map from blocking tasks to the tasks they block
	blockedBy map[task.Id][]task.Id // map from blocked tasks to the tasks they are blocked by
}

// NewTaskTree creates a new, empty TaskTree.
func NewTaskTree() *TaskTree {
	return &TaskTree{
		tasks:     make(map[task.Id]task.Task),
		roots:     make([]task.Id, 0),
		subtasks:  make(map[task.Id][]task.Id),
		subtaskOf: make(map[task.Id]task.Id),
		blocks:    make(map[task.Id][]task.Id),
		blockedBy: make(map[task.Id][]task.Id),
	}
}

func (tree *TaskTree) assertTaskExists(id task.Id) error {
	if _, exists := tree.tasks[id]; !exists {
		return fmt.Errorf("task %v does not exist", id)
	}

	return nil
}

func (tree *TaskTree) idsToTasks(ids []task.Id) []task.Task {
	return util.Map(ids, func(id task.Id) task.Task {
		return tree.tasks[id]
	})
}

// AddTask adds a Task object to the TaskTree.
func (tree *TaskTree) AddTask(task task.Task) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if _, exists := tree.tasks[task.Id]; exists {
		return fmt.Errorf("task with id %v already exists", task.Id)
	}

	tree.tasks[task.Id] = task
	tree.roots = append(tree.roots, task.Id)
	return nil
}

// GetTask gets a Task object in the TaskTree by its id.
func (tree *TaskTree) GetTask(id task.Id) (task task.Task, exists bool) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()
	task, exists = tree.tasks[id]
	return
}

// DeleteTask deletes a task from the tree by id.
func (tree *TaskTree) DeleteTask(id task.Id) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if err := tree.assertTaskExists(id); err != nil {
		return err
	}

	parentId, isSubtask := tree.subtaskOf[id]
	if !isSubtask {
		tree.roots = util.Remove(tree.roots, id)
	} else {
		tree.subtasks[parentId] = util.Remove(tree.subtasks[parentId], id)
		delete(tree.subtaskOf, id)
	}

	delete(tree.tasks, id)
	for _, subtaskId := range tree.subtasks[id] {
		delete(tree.subtaskOf, subtaskId)
		tree.roots = append(tree.roots, subtaskId)
	}
	delete(tree.subtasks, id)
	return nil
}

// UpdateTask replaces a task in the tree with an updated version.
func (tree *TaskTree) UpdateTask(task task.Task) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if _, exists := tree.tasks[task.Id]; !exists {
		return fmt.Errorf("task with id %v does not exist", task.Id)
	}

	tree.tasks[task.Id] = task
	return nil
}

// GetRootTasks returns the root tasks (i.e. tasks that aren't subtasks/don't have parents).
func (tree *TaskTree) GetRootTasks() []task.Task {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()
	return util.Map(tree.roots, func(id task.Id) task.Task { return tree.tasks[id] })
}
