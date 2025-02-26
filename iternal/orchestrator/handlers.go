package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/OnYyon/GoroutineRPNServer/iternal/parser"
	"github.com/google/uuid"
)

// For global data storage
func NewAPI() *API {
	return &API{
		Expressions: make(map[string]Expression),
		Tasks:       make(map[string][]Task),
	}
}

func getID() string {
	return uuid.New().String()
}

func createTasks(rpn []string, expID string) []Task {
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

			task := Task{ID: getID(), Arg1: a, Arg2: b, Operation: v, Status: StatusNew, ExpressionID: expID}
			tasks = append(tasks, task)

			stack = append(stack, task.ID)
		}
	}
	if len(stack) != 1 {
		return nil
	}
	return tasks
}

func (a *API) AddNewExpression(w http.ResponseWriter, r *http.Request) {
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
	tasks := createTasks(rpn, expression.ID)
	fmt.Println(tasks)
	a.mu.Lock()
	a.Tasks[expression.ID] = tasks
	a.mu.Unlock()
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"id": expression.ID})
}

func (a *API) GetSliceOfExpressions(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *API) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	expression, ok := a.Expressions[id]
	if !ok {
		http.Error(w, "Incorrect ID", 404)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expression)
}
