package orchestrator

import "sync"

type API struct {
	Expressions map[string]Expression
	Tasks       map[string]Task
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
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	Result        float64 `json:"result"`
	Status        string  `json:"status"`
	OperationTime int64   `json:"operation_time"`
	ExpressionID  string  `json:"expression_id"`
}

const (
	StatusNew        = "new"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)
