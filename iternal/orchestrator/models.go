package orchestrator

import (
	"sync"
	"time"

	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
)

type API struct {
	Expressions map[string]Expression
	Tasks       map[string][]Task
	queque      chan Task
	cfg         *config.Config
	mu          sync.Mutex
}

type Expression struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
	Tasks  []Task  `json:"-"`
	Input  string  `json:"-"`
}

type Task struct {
	ID            string        `json:"id"`
	Arg1          string        `json:"arg1"`
	Arg2          string        `json:"arg2"`
	Operation     string        `json:"operation"`
	Result        float64       `json:"result"`
	Status        string        `json:"-"`
	OperationTime time.Duration `json:"operation_time"`
	ExpressionID  string        `json:"-"`
	Error         error
}

const (
	StatusNew        = "new"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)
