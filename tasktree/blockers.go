package tasktree

import (
	"fmt"
	"github.com/carreter/tasktree-go/task"
	"github.com/carreter/tasktree-go/util"
)

// MarkBlocker marks one task (blocker) as a prerequisite for another task (blocked).
func (tree *TaskTree) MarkBlocker(blockerId task.Id, blockedId task.Id) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	// TODO: Add detection of circular blocking chains

	if err := tree.assertTaskExists(blockerId); err != nil {
		return err
	}
	if err := tree.assertTaskExists(blockedId); err != nil {
		return err
	}

	tree.blocks[blockerId] = append(tree.blocks[blockerId], blockedId)
	tree.blockedBy[blockedId] = append(tree.blockedBy[blockedId], blockerId)
	return nil
}

// UnmarkBlocker marks one task (blocker) as a no longer being a prerequisite for another task (blocked).
func (tree *TaskTree) UnmarkBlocker(blockerId task.Id, blockedId task.Id) error {
	tree.rwMu.Lock()
	defer tree.rwMu.Unlock()

	if err := tree.assertTaskExists(blockerId); err != nil {
		return err
	}
	if err := tree.assertTaskExists(blockedId); err != nil {
		return err
	}

	if blocks, exists := tree.blocks[blockerId]; exists {
		tree.blocks[blockerId] = util.Remove(blocks, blockedId)
	}
	if blockedBy, exists := tree.blockedBy[blockedId]; exists {
		tree.blocks[blockerId] = util.Remove(blockedBy, blockerId)
	}

	return nil
}

// GetDirectBlockers gets tasks that are directly blocking the specific task.
func (tree *TaskTree) GetDirectBlockers(id task.Id) ([]task.Task, error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	if err := tree.assertTaskExists(id); err != nil {
		return nil, err
	}

	blockerIds, exists := tree.blockedBy[id]
	if !exists {
		return nil, nil
	}

	return tree.idsToTasks(blockerIds), nil
}

// GetAllBlockers gets all tasks that either:
//   - directly block a task
//   - block a task's ancestors
func (tree *TaskTree) GetAllBlockers(id task.Id) ([]task.Task, error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	// TODO: Ensure blockers are returned in a logical order (toposort?)

	if err := tree.assertTaskExists(id); err != nil {
		return nil, err
	}

	seenBlockers := make(map[task.Id]struct{})

	blockers, err := tree.GetDirectBlockers(id)
	if err != nil {
		return nil, fmt.Errorf("could not get direct blockers of task %v", id)
	}
	for _, blocker := range blockers {
		seenBlockers[blocker.Id] = struct{}{}
	}

	ancestors, err := tree.GetAncestorTasks(id)
	if err != nil {
		return nil, err
	}

	for _, ancestor := range ancestors {
		ancestorBlockers, err := tree.GetDirectBlockers(ancestor.Id)
		if err != nil {
			return nil, err
		}

		for _, ancestorBlocker := range ancestorBlockers {
			if _, exists := seenBlockers[ancestorBlocker.Id]; !exists {
				blockers = append(blockers, ancestorBlocker)
				seenBlockers[ancestorBlocker.Id] = struct{}{}
			}
		}
	}

	return blockers, nil
}

// IsBlocked checks if a task or any of its parent tasks are blocked.
func (tree *TaskTree) IsBlocked(id task.Id) (bool, error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	blockers, err := tree.GetAllBlockers(id)
	if err != nil {
		return false, err
	}

	return len(blockers) != 0, nil
}
