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
		Expressions: make(map[string]*Expression),
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

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func (a *API) createTasks(rpn []string, expID string) ([]string, []Task) {
	stack := []string{}
	tasks := []Task{}
	for _, v := range rpn {
		if isOperator(v) {
			if len(stack) < 2 {
				stack = append(stack, v)
				continue
			}
			if isNumber(stack[len(stack)-1]) && isNumber(stack[len(stack)-2]) {
				task := Task{
					ID:            getID(),
					Arg1:          stack[len(stack)-2],
					Arg2:          stack[len(stack)-1],
					Operation:     v,
					Status:        StatusNew,
					OperationTime: getTimeOp(v, a.cfg),
					ExpressionID:  expID,
				}
				tasks = append(tasks, task)
				stack = stack[:len(stack)-2]
				stack = append(stack, task.ID)
			} else {
				stack = append(stack, v)
			}
		} else if isNumber(v) {
			stack = append(stack, v)
		} else {
			a.mu.Lock()
			pos := 0
			for i, j := range a.Tasks[expID] {
				if j.ID == v {
					pos = i
				}
			}
			stack = append(stack, fmt.Sprint(a.Tasks[expID][pos].Result))
			a.mu.Unlock()
		}
	}
	return stack, tasks
}

func (a *API) continueExpressionCalculation(expID string) {
	new_rpn, tasks := a.createTasks(a.rpnCurrent, expID)
	if len(tasks) == 0 {
		a.Expressions[expID].Result = a.Tasks[expID][len(a.Tasks[expID])-1].Result
		//fmt.Printf("Final result for expression %s: %v\n", expID, a.Tasks[expID][len(a.Tasks[expID])-1].Result)
		return
	}

	a.mu.Lock()
	a.Tasks[expID] = append(a.Tasks[expID], tasks...)
	a.rpnCurrent = new_rpn
	a.mu.Unlock()

	go func() {
		for _, v := range tasks {
			a.queque <- v
		}
	}()
}

func (a *API) calculateExpression(exp *Expression) {
	rpn, err := parser.ParserToRPN(exp.Input)
	if err != nil {
		return
	}
	new_rpn, tasks := a.createTasks(rpn, exp.ID)
	a.mu.Lock()
	a.Tasks[exp.ID] = append(a.Tasks[exp.ID], tasks...)
	a.rpnCurrent = new_rpn
	a.mu.Unlock()
	go func() {
		for _, v := range tasks {
			a.queque <- v
		}
	}()
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
	a.Expressions[expression.ID] = &expression
	a.mu.Unlock()
	a.calculateExpression(&expression)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"id": expression.ID})
}

func (a *API) GetSliceOfExpressions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	expressionsSlice := make([]Expression, 0, len(a.Expressions))
	for _, expr := range a.Expressions {
		expressionsSlice = append(expressionsSlice, *expr)
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
	json.NewEncoder(w).Encode(&task)
}

func (a *API) GetPostResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	for exprID, tasks := range a.Tasks {
		for i, task := range tasks {
			if task.ID == result.ID {
				a.Tasks[exprID][i].Result = result.Result
				a.Tasks[exprID][i].Status = StatusCompleted
				// fmt.Printf("Updated task %s with result: %f\n", result.ID, result.Result)

				allCompleted := true
				for _, t := range a.Tasks[exprID] {
					if t.Status != StatusCompleted {
						allCompleted = false
						break
					}
				}

				if allCompleted {
					// fmt.Println("All tasks completed. Processing next steps.")
					a.Expressions[exprID] = &Expression{
						ID:     exprID,
						Status: StatusCompleted,
						Input:  a.Expressions[exprID].Input,
					}
					go a.continueExpressionCalculation(exprID)
				}
				w.WriteHeader(200)
				return
			}
		}
	}
	http.Error(w, "Task not found", 404)
}
