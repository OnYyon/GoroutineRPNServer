package orchestrator

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
	"github.com/OnYyon/GoroutineRPNServer/iternal/parser"
	"github.com/google/uuid"
)

func NewAPI() *API {
	return &API{
		Expressions: make(map[string]*Expression),
		Tasks:       make(map[string][]Task),
		rpnCurrent:  make(map[string][]string),
		repeats:     make(map[string]int),
		cfg:         config.NewConfig(),
		queque:      make(chan Task),
	}
}

func (a *API) getTaskFromChan() (Task, bool) {
	a.muTasks.Lock()
	defer a.muTasks.Unlock()
	select {
	case task := <-a.queque:
		return task, true
	default:
		return Task{}, false
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

func CheckExpression(exp string) bool {
	// Проверка на разрешенные символы без пробелов и правильное размещение операторов
	re := regexp.MustCompile(`^[0-9\+\-\*/\(\)]*$`) // Разрешаем только цифры, операторы и скобки
	if !re.MatchString(exp) {
		return false
	}

	// Проверка на правильную последовательность чисел, операторов и скобок
	// Пример: (7+1)/(2+2)*4, (32-((4+12)*2))-1
	re2 := regexp.MustCompile(`^(\d+|\([\d\+\-\*/\(\)]+\))([\+\-\*/]\d+|\([\d\+\-\*/\(\)]+\))*$`)
	return re2.MatchString(exp)
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
			a.muTasks.Lock()
			pos := 0
			for i, j := range a.Tasks[expID] {
				if j.ID == v {
					pos = i
				}
			}
			stack = append(stack, fmt.Sprint(a.Tasks[expID][pos].Result))
			a.muTasks.Unlock()
		}
	}
	return stack, tasks
}

func (a *API) calculateExpression(exp *Expression) {
	a.Expressions[exp.ID].Status = StatusInProgress
	rpn, err := parser.ParserToRPN(exp.Input)
	if err != nil {
		a.Expressions[exp.ID].Status = StatusFailed
		return
	}
	new_rpn, tasks := a.createTasks(rpn, exp.ID)
	a.muTasks.Lock()
	a.Tasks[exp.ID] = tasks
	a.rpnCurrent[exp.ID] = new_rpn
	a.muTasks.Unlock()
	go func() {
		for _, v := range tasks {
			a.queque <- v
		}
	}()
}

func (a *API) continueExpressionCalculation(expID string) {
	new_rpn, tasks := a.createTasks(a.rpnCurrent[expID], expID)
	if len(tasks) == 0 {
		a.Expressions[expID].Result = a.Tasks[expID][len(a.Tasks[expID])-1].Result
		// For debugging
		// fmt.Printf("Final result for expression %s: %v\n", expID, a.Tasks[expID][len(a.Tasks[expID])-1].Result)
		return
	}

	a.muTasks.Lock()
	a.Tasks[expID] = append(a.Tasks[expID], tasks...)
	a.rpnCurrent[expID] = new_rpn

	go func() {
		for _, v := range tasks {
			a.queque <- v
		}
	}()
	a.muTasks.Unlock()
}
