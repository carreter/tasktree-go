package app

import (
	"github.com/carreter/tasktree-go/pkg/tasktree"
	"sync"
)

// Context provides global context for the app.
type Context struct {
	mu       sync.Mutex
	taskTree *tasktree.TaskTree
}

func NewContext(taskTree *tasktree.TaskTree) *Context {
	return &Context{taskTree: taskTree}
}

func (ctx *Context) TaskTree() *tasktree.TaskTree {
	return ctx.taskTree
}

func (ctx *Context) SetTaskTree(taskTree *tasktree.TaskTree) {
	ctx.taskTree = taskTree
}
