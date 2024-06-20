package tasktree

import (
	"fmt"
	"github.com/google/uuid"
	"hobby-tracker/pkg/util"
)

// MarkSubtask marks one task (subtask) as a subtask of another (parent).
func (tree *TaskTree) MarkSubtask(parentId uuid.UUID, subtaskId uuid.UUID) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if err := tree.assertTaskExists(parentId); err != nil {
		return err
	}
	if err := tree.assertTaskExists(subtaskId); err != nil {
		return err
	}

	if existingParentId, exists := tree.subtaskOf[subtaskId]; exists {
		return fmt.Errorf("task %v is already a subtask of %v", subtaskId, existingParentId)
	}

	// TODO: Add subtask cycle detection.

	tree.subtasks[parentId] = append(tree.subtasks[parentId], subtaskId)
	tree.subtaskOf[subtaskId] = parentId
	return nil
}

// UnmarkSubtask marks a task as an independent task rather than a subtask.
// Does not error if task was already an independent task.
func (tree *TaskTree) UnmarkSubtask(subtaskId uuid.UUID) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if err := tree.assertTaskExists(subtaskId); err != nil {
		return err
	}

	parentId, exists := tree.subtaskOf[subtaskId]
	if !exists {
		return nil
	}

	delete(tree.subtaskOf, subtaskId)
	tree.subtasks[parentId] = util.Remove(tree.subtasks[parentId], subtaskId)
	return nil
}

// GetDirectSubtasksOf gets the direct children of a Task.
func (tree *TaskTree) GetDirectSubtasksOf(parentId uuid.UUID) ([]Task, error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	if err := tree.assertTaskExists(parentId); err != nil {
		return nil, err
	}

	subtaskIds, exists := tree.subtasks[parentId]
	if !exists {
		return nil, nil
	}

	return tree.idsToTasks(subtaskIds), nil
}

// GetParentTask gets the parent of a given task if it exists.
func (tree *TaskTree) GetParentTask(id uuid.UUID) (parent Task, exists bool, err error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	if err := tree.assertTaskExists(id); err != nil {
		return Task{}, false, err
	}

	parentId, exists := tree.subtaskOf[id]
	if !exists {
		return Task{}, false, nil
	}

	return tree.tasks[parentId], true, nil
}

// IsSubtask determines whether a Task is a subtask (i.e. has a parent Task).
func (tree *TaskTree) IsSubtask(id uuid.UUID) (bool, error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	if err := tree.assertTaskExists(id); err != nil {
		return false, err
	}

	_, exists, err := tree.GetParentTask(id)
	return exists, err
}

// GetAncestorTasks returns the ancestors of a task in order
// (parent task, grandparent task, great grandparent task, etc.)
func (tree *TaskTree) GetAncestorTasks(id uuid.UUID) ([]Task, error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	if err := tree.assertTaskExists(id); err != nil {
		return nil, err
	}

	res := make([]Task, 0)
	currId := id
	for {
		ancestor, exists, err := tree.GetParentTask(currId)
		if err != nil {
			return nil, err
		} else if !exists {
			return res, nil
		}

		res = append(res, ancestor)
		currId = ancestor.Id
	}
}
