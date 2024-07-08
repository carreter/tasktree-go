package command

import (
	"fmt"
	"github.com/carreter/tasktree-go/app"
	"github.com/carreter/tasktree-go/pkg/task"
)

type DeleteCommand struct {
}

func (c DeleteCommand) Run(ctx *app.Context, args ...string) (string, string) {
	if len(args) != 2 {
		return "", fmt.Sprintf("inccorect number of arguments, usage: %v", c.Usage())
	}

	taskId := task.Id(args[1])
	err := ctx.TaskTree().DeleteTask(taskId)
	if err != nil {
		return "", fmt.Sprintf("failed to delete task: %v", err)
	}

	return fmt.Sprintf("deleted task %v", taskId), ""
}

func (c DeleteCommand) Usage() string {
	return "add <task name> [<task description>]"
}

func (c DeleteCommand) Name() string {
	return "add"
}
