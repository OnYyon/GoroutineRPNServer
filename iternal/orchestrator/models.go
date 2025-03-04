package orchestrator

import (
	"sync"
	"time"

	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
)

type API struct {
	Expressions map[string]*Expression
	Tasks       map[string][]Task
	rpnCurrent  map[string][]string
	repeats     map[string]int
	queque      chan Task
	cfg         *config.Config
	muTasks     sync.RWMutex
	muExpr      sync.RWMutex
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

type Res struct {
	ID      string
	Result  float64
	Timeout bool
	Errors  error
}

const (
	StatusNew        = "new"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)
