package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
	"github.com/OnYyon/GoroutineRPNServer/iternal/parser"
	"github.com/google/uuid"
)

// For global data storage
func NewAPI() *API {
	return &API{
		Expressions: make(map[string]Expression),
		Tasks:       make(map[string][]Task),
		cfg:         config.NewConfig(),
		queque:      make(chan Task),
	}
}

func getID() string {
	return uuid.New().String()
}

func getTimeOp(operation string, cfg *config.Config) time.Duration {
	if operation == "+" {
		return cfg.TimeAdd
	} else if operation == "-" {
		return cfg.TimeSubtraction
	} else if operation == "*" {
		return cfg.TimeMultiply
	} else {
		return cfg.TimeDivision
	}
}

func createTasks(rpn []string, expID string, cfg *config.Config) []Task {
	var stack []string
	var tasks []Task

	for _, v := range rpn {
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			stack = append(stack, v)
		} else {
			if len(stack) < 2 {
				return nil
			}
			a, b := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			task := Task{
				ID:            getID(),
				Arg1:          a,
				Arg2:          b,
				Operation:     v,
				Status:        StatusNew,
				OperationTime: getTimeOp(v, cfg), // NOTE: I think it can be done better.
				ExpressionID:  expID,
			}
			tasks = append(tasks, task)

			stack = append(stack, task.ID)
		}
	}
	if len(stack) != 1 {
		return nil
	}
	return tasks
}

func (a *API) getTaskFromChan() (Task, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	select {
	case task := <-a.queque:
		return task, true
	default:
		return Task{}, false
	}
}

func (a *API) AddNewExpression(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Oppps something went wrong", 500)
		return
	}

	expression := Expression{
		ID:     getID(),
		Status: StatusNew,
		Input:  request.Expression,
	}
	a.mu.Lock()
	a.Expressions[expression.ID] = expression
	a.mu.Unlock()
	rpn, err := parser.ParserToRPN(expression.Input)
	if err != nil {
		http.Error(w, "Oppps something went wrong", 500)
		return
	}
	tasks := createTasks(rpn, expression.ID, a.cfg)
	fmt.Println(tasks)
	a.mu.Lock()
	a.Tasks[expression.ID] = tasks

	// For tests
	go func() {
		for _, v := range tasks {
			a.queque <- v
		}
		// FIXME:
		// close(a.queque)
	}()
	a.mu.Unlock()
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"id": expression.ID})
}

func (a *API) GetSliceOfExpressions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	expressionsSlice := make([]Expression, 0, len(a.Expressions))
	for _, expr := range a.Expressions {
		expressionsSlice = append(expressionsSlice, expr)
	}
	response := struct {
		Expressions []Expression `json:"expressions"`
	}{
		Expressions: expressionsSlice,
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

func (a *API) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: Varible in url. Change method maybe
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	expression, ok := a.Expressions[id]
	if !ok {
		http.Error(w, "Incorrect ID", 404)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(expression)
}

// Iternal part of handlers
func (a *API) GetTasksToAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	task, ok := a.getTaskFromChan()
	if !ok {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(task)
}
