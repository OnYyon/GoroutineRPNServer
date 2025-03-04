package agent

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
)

func Worker(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for {
		task, err := FetchTask()
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		result, err := evaluateExpression(task)
		if err != nil {
			continue
		}
		err = SendResult(Res{ID: task.ID, Result: result})
		fmt.Printf("sending %v, goroutine: %v\n", task.ID, id)
		if err != nil {
			fmt.Printf("Worker %d: failed to send result", err)
		}
	}
}

func evaluateExpression(task *orchestrator.Task) (float64, error) {
	a, err := strconv.ParseFloat(task.Arg1, 64)
	if err != nil {
		return 0, err
	}
	b, err := strconv.ParseFloat(task.Arg2, 64)
	if err != nil {
		return 0, err
	}
	var result float64
	switch task.Operation {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("zero dividion")
		}
		result = a / b
	}
	return result, nil
}
