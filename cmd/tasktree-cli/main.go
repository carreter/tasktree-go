package main

import (
	"fmt"
	"github.com/carreter/tasktree-go/app/models"
	"github.com/carreter/tasktree-go/pkg/tasktree"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	model := models.NewModel(tasktree.NewTaskTree())
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Printf("fatal error: %v", err)
		os.Exit(1)
	}
}
