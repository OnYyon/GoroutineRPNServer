package agent

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
)

func Worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, have := GetTask()
		if !have {
			time.Sleep(1 * time.Second)
			continue
		}
		err := calculateTask(&task)
		if err != nil {
			task.Error = err
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("Received task: %+v", task)
	}
}

func calculateTask(task *orchestrator.Task) error {
	a, err := strconv.ParseFloat(task.Arg1, 64)
	if err != nil {
		return err
	}
	b, err := strconv.ParseFloat(task.Arg2, 64)
	if err != nil {
		return err
	}
	switch task.Operation {
	case "+":
		task.Result = a + b
		fmt.Println(task.Result)
	case "-":
		task.Result = a - b
	case "*":
		task.Result = a * b
		fmt.Println(task.Result)
	case "/":
		if b == 0 {
			return fmt.Errorf("division by zero")
		}
		task.Result = a / b
	}
	return nil
}
