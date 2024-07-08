package command

import (
	"fmt"
	"github.com/carreter/tasktree-go/app"
)

type HelpCommand struct {
	Commands *map[string]Command
}

func (c HelpCommand) Run(ctx *app.Context, args ...string) (string, string) {
	if len(args) != 2 {
		return "", fmt.Sprintf("inccorect number of arguments, usage: %v", c.Usage())
	}

	cmd, exists := (*c.Commands)[args[1]]
	if !exists {
		return "", fmt.Sprintf("command not found: %s", args[1])
	}

	return fmt.Sprintf("usage: %v", cmd.Usage()), ""
}

func (c HelpCommand) Usage() string {
	return "help <command>"
}

func (c HelpCommand) Name() string {
	return "help"
}
