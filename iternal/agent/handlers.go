package agent

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
)

func FetchTask() (*orchestrator.Task, error) {
	resp, err := http.Get("http://127.0.0.1:8080/iternal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var task orchestrator.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func SendResult(result Res) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://127.0.0.1:8080/iternal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
