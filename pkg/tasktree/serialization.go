package tasktree

import (
	"bytes"
	"encoding/gob"
	"github.com/carreter/tasktree-go/pkg/task"
)

// GobEncode allows for gob encoding of a TaskTree.
// A custom implementation is necessary here because the TaskTree
// struct fields are private.
func (tree *TaskTree) GobEncode() ([]byte, error) {
	w := &bytes.Buffer{}
	encoder := gob.NewEncoder(w)

	// We only need to encode the tree.tasks, tree.subtasks, and tree.blocks
	// maps as we can reconstruct tree.subtaskOf and tree.blockedBy from these.
	err := encoder.Encode(tree.tasks)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(tree.subtasks)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(tree.blocks)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

// GobDecode allows for gob decoding of a TaskTree.
// A custom implementation is necessary here because the TaskTree
// struct fields are private.
func (tree *TaskTree) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)

	err := decoder.Decode(&tree.tasks)
	if err != nil {
		return err
	}
	err = decoder.Decode(&tree.subtasks)
	if err != nil {
		return err
	}
	err = decoder.Decode(&tree.blocks)
	if err != nil {
		return err
	}

	tree.rehydrate()

	return nil
}

// rehydrate reconstructs the tree.subtaskOf and tree.blockedBy maps.
func (tree *TaskTree) rehydrate() {
	for parentId, subtaskIds := range tree.subtasks {
		for _, subtaskId := range subtaskIds {
			tree.subtaskOf[subtaskId] = parentId
		}
	}

	for blockerId, blockedIds := range tree.blocks {
		for _, blockedId := range blockedIds {
			blockedBy, exists := tree.blockedBy[blockedId]
			if !exists {
				blockedBy = make([]task.Id, 0)
			}
			tree.blockedBy[blockedId] = append(blockedBy, blockerId)
		}
	}
}
