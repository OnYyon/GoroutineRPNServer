package orchestrator

// TODO: logger

import (
	"encoding/json"
	"net/http"
	"strings"
)

// For global data storage
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
	a.muExpr.Lock()
	a.Expressions[expression.ID] = &expression
	a.muExpr.Unlock()
	go a.calculateExpression(&expression)
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
		ID      string  `json:"id"`
		Result  float64 `json:"result"`
		Timeout bool    `json:"timeout"`
		Errors  error   `json:"errors"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	a.muTasks.Lock()
	defer a.muTasks.Unlock()
	for exprID, tasks := range a.Tasks {
		for i, task := range tasks {
			if result.Errors != nil {
				a.Tasks[exprID][i].Status = StatusFailed
				return
			} else if result.Timeout {
				if a.repeats[result.ID] >= 5 {
					a.Tasks[exprID][i].Status = StatusFailed
				}
				a.repeats[result.ID]++
				go func() {
					a.queque <- task
				}()
				return
			}
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
