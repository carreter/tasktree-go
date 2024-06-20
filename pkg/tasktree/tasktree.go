package tasktree

import (
	"fmt"
	"github.com/google/uuid"
	"hobby-tracker/pkg/util"
	"sync"
)

// A TaskTree organizes the relationships between tasks. Thread-safe.
type TaskTree struct {
	rwMu sync.RWMutex

	tasks map[uuid.UUID]Task

	subtasks  map[uuid.UUID][]uuid.UUID
	subtaskOf map[uuid.UUID]uuid.UUID

	blocks    map[uuid.UUID][]uuid.UUID
	blockedBy map[uuid.UUID][]uuid.UUID
}

// NewTaskTree creates a new, empty TaskTree.
func NewTaskTree() *TaskTree {
	return &TaskTree{
		tasks:    make(map[uuid.UUID]Task),
		subtasks: make(map[uuid.UUID][]uuid.UUID),
		blocks:   make(map[uuid.UUID][]uuid.UUID),
	}
}

func (tree *TaskTree) assertTaskExists(id uuid.UUID) error {
	if _, exists := tree.tasks[id]; !exists {
		return fmt.Errorf("task %v does not exist", id)
	}

	return nil
}

func (tree *TaskTree) idsToTasks(ids []uuid.UUID) []Task {
	return util.Map(ids, func(id uuid.UUID) Task {
		return tree.tasks[id]
	})
}

// AddTask adds a Task object to the TaskTree.
func (tree *TaskTree) AddTask(task Task) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if _, exists := tree.tasks[task.Id]; exists {
		return fmt.Errorf("task with id %v already exists", task.Id)
	}

	tree.tasks[task.Id] = task
	return nil
}

// GetTask gets a Task object in the TaskTree by its id.
func (tree *TaskTree) GetTask(id uuid.UUID) (task Task, exists bool) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()
	task, exists = tree.tasks[id]
	return
}

// DeleteTask deletes a task from the tree by id.
func (tree *TaskTree) DeleteTask(id uuid.UUID) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if err := tree.assertTaskExists(id); err != nil {
		return err
	}

	delete(tree.tasks, id)
	return nil
}

// UpdateTask replaces a task in the tree with an updated version.
func (tree *TaskTree) UpdateTask(task Task) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if _, exists := tree.tasks[task.Id]; !exists {
		return fmt.Errorf("task with id %v does not exist", task.Id)
	}

	tree.tasks[task.Id] = task
	return nil
}
