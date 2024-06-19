package tasktree

import (
	"fmt"
	"github.com/google/uuid"
	"hobby-tracker/pkg/util"
)

// MarkBlocker marks one task (blocker) as a prerequisite for another task (blocked).
func (tree *TaskTree) MarkBlocker(blockerId uuid.UUID, blockedId uuid.UUID) error {
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
func (tree *TaskTree) UnmarkBlocker(blockerId uuid.UUID, blockedId uuid.UUID) error {
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
func (tree *TaskTree) GetDirectBlockers(id uuid.UUID) ([]Task, error) {
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
func (tree *TaskTree) GetAllBlockers(id uuid.UUID) ([]Task, error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	// TODO: Ensure blockers are returned in a logical order (toposort?)

	if err := tree.assertTaskExists(id); err != nil {
		return nil, err
	}

	seenBlockers := make(map[uuid.UUID]struct{})

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
func (tree *TaskTree) IsBlocked(id uuid.UUID) (bool, error) {
	tree.rwMu.RLock()
	defer tree.rwMu.RUnlock()

	blockers, err := tree.GetAllBlockers(id)
	if err != nil {
		return false, err
	}

	return len(blockers) != 0, nil
}
