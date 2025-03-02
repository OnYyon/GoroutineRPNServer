package agent

import (
	"encoding/json"
	"net/http"

	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
)

func GetTask() (orchestrator.Task, bool) {
	resp, err := http.Get("http://127.0.0.1:8080/iternal/task")
	if err != nil {
		return orchestrator.Task{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return orchestrator.Task{}, false
	}
	var task orchestrator.Task
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return orchestrator.Task{}, false
	}
	return task, true
}
