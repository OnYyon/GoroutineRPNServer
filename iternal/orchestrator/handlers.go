package orchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// For global data storage
func NewAPI() *API {
	return &API{
		Expressions: make(map[string]Expression),
		Tasks:       make(map[string]Task),
	}
}

func newID() string {
	return uuid.New().String()
}

func (a *API) addNewExpression(w http.ResponseWriter, r *http.Response) {
	var request struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expression := Expression{
		ID:     newID(),
		Status: StatusNew,
		Input:  request.Expression,
	}
	a.mu.Lock()
	a.Expressions[expression.ID] = expression
	a.mu.Unlock()
}
