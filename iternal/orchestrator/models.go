package orchestrator

import (
	"sync"
	"time"

	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
)

type API struct {
	Expressions map[string]*Expression
	Tasks       map[string][]Task
	queque      chan Task
	rpnCurrent  map[string][]string
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

type Result struct {
	ID     string
	Result float64
}

const (
	StatusNew        = "new"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

/*
curl -X POST http://127.0.0.1:8080/api/v1/calculate \
     -H "Content-Type: application/json" \
     -d '{"expression": "(2+2)*2"}' &

curl -X POST http://127.0.0.1:8080/api/v1/calculate \
     -H "Content-Type: application/json" \
     -d '{"expression": "2+2*2"}' &
*/
