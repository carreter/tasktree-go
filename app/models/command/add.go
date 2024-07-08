package commands

import (
	"fmt"
	"github.com/carreter/tasktree-go/app"
	"github.com/carreter/tasktree-go/pkg/task"
	"strings"
)

type AddCommand struct {
}

func (c AddCommand) Run(ctx *app.Context, args ...string) (string, string) {
	if len(args) < 2 {
		return "", fmt.Sprintf("inccorect number of arguments, usage: %v", c.Usage())
	}

	newTask := task.Task{Id: task.Id(args[1]), Name: args[1]}
	if len(args) > 2 {
		newTask.Description = strings.Trim(strings.Join(args[2:], " "), "\"")
	}

	err := ctx.TaskTree().AddTask(newTask)
	if err != nil {
		return "", fmt.Sprintf("failed to add task: %v", err)
	}

	return fmt.Sprintf("added task with ID %v", newTask.Id), ""
}

func (c AddCommand) Usage() string {
	return "add <task name> [<task description>]"
}

func (c AddCommand) Name() string {
	return "add"
}
