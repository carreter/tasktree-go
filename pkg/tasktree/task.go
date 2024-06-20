// Package tasktree defines a tree of tasks.
package tasktree

import (
	"github.com/google/uuid"
	"time"
)

// A Tag is a searchable piece of information about a Task.
type Tag string

// A Priority represents how important it is for a Task to be completed.
type Priority byte

const (
	// Default priority.
	Default Priority = iota
	// Urgent priority.
	Urgent
	// High priority.
	High
	// Normal priority.
	Normal
	// Low priority.
	Low
)

// A Task represents an individual task.
type Task struct {
	Id            uuid.UUID
	Name          string
	Description   string
	EstimatedTime time.Duration
	TimeInvested  time.Duration
	Completed     bool
	Tags          []Tag
	Priority      Priority
}
