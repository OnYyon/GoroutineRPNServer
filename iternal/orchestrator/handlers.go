package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
	"github.com/OnYyon/GoroutineRPNServer/iternal/parser"
	"github.com/google/uuid"
)

var Wg sync.WaitGroup

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

// LEGACY:
/*
func createTasks(rpn []string, expID string, cfg *config.Config) []Task {
	var stack []*float64
	var tasks []Task

	for _, v := range rpn {
		if op, err := strconv.ParseFloat(v, 64); err == nil {
			stack = append(stack, &op)
		} else {
			if len(stack) < 2 {
				return nil
			}
			a, b := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			task := Task{
				ID:            getID(),
				Arg1:          *a,
				Arg2:          *b,
				Operation:     v,
				Status:        StatusNew,
				OperationTime: getTimeOp(v, cfg), // NOTE: I think it can be done better.
				ExpressionID:  expID,
			}
			tasks = append(tasks, task)

			stack = append(stack, &task.Result)
		}
	}
	if len(stack) != 1 {
		return nil
	}
	for _, v := range tasks {
		fmt.Println(v.ID, v.Arg1, v.Arg1, v.Result, v.Status)
	}
	return tasks
}
*/

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func (a *API) createTasks(rpn []string, expID string, pos int) ([]string, []Task) {
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
			stack = append(stack, fmt.Sprint(a.Tasks[expID][pos].Result))
			a.mu.Unlock()
			fmt.Println(stack)
			pos++
		}
	}

	return stack, tasks
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
	fmt.Println(rpn)
	if err != nil {
		http.Error(w, "Oppps something went wrong", 500)
		return
	}
	pos := 0
	// TODO: Пересмотреть структуру добавления
	new_rpn, tasks := a.createTasks(rpn, expression.ID, pos)
	for len(new_rpn) != 1 {
		a.mu.Lock()
		a.Tasks[expression.ID] = append(a.Tasks[expression.ID], tasks...)
		go func() {
			for _, v := range tasks {
				Wg.Add(1)
				a.queque <- v
			}
			Wg.Wait()
		}()
		a.mu.Unlock()
		new_rpn, tasks = a.createTasks(new_rpn, expression.ID, pos)
	}
	fmt.Println(a.Tasks[expression.ID])
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

func (a *API) GetPostResult(w http.ResponseWriter, r *http.Request) {

}
