package command

import (
	"fmt"
	"github.com/carreter/tasktree-go/app"
	"github.com/carreter/tasktree-go/pkg/task"
	"strings"
)

type AddSubtaskCommand struct {
}

func (c AddSubtaskCommand) Run(ctx *app.Context, args ...string) (string, string) {
	if len(args) < 3 {
		return "", fmt.Sprintf("inccorect number of arguments, usage: %v", c.Usage())
	}

	newTask := task.Task{Id: task.Id(args[2]), Name: args[2]}
	if len(args) > 3 {
		newTask.Description = strings.Trim(strings.Join(args[3:], " "), "\"")
	}

	err := ctx.TaskTree().AddTask(newTask)
	if err != nil {
		return "", fmt.Sprintf("failed to add task: %v", err)
	}
	parentTaskId := task.Id(args[1])
	err = ctx.TaskTree().MarkSubtask(parentTaskId, newTask.Id)
	if err != nil {
		return "", fmt.Sprintf("failed to add subtask: %v", err)
	}

	return fmt.Sprintf("added task %v as subtask of task %v", newTask.Id, parentTaskId), ""
}

func (c AddSubtaskCommand) Usage() string {
	return "add-subtask <task name> [<task description>]"
}

func (c AddSubtaskCommand) Name() string {
	return "add-subtask"
}
