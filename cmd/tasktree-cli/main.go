package main

import (
	"fmt"
	"github.com/carreter/tasktree-go/app/views"
	"github.com/carreter/tasktree-go/pkg/tasktree"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	model := views.NewModel(tasktree.NewTaskTree())
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Printf("fatal error: %v", err)
		os.Exit(1)
	}
}
