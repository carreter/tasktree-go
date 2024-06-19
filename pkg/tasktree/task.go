package tasktree

import (
	"github.com/google/uuid"
	"time"
)

type Tag string

// A Task represents an individual task.
type Task struct {
	Id            uuid.UUID
	Name          string
	Description   string
	EstimatedTime time.Duration
	TimeInvested  time.Duration
	Completed     bool
	Tags          []Tag
}
