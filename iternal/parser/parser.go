package parser

import (
	"fmt"
	"strings"
	"unicode"
)

func ParserToRPN(pattern string) ([]string, error) {
	var lst []string
	var stack, ouput_stack []string
	operators := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}
	pattern = strings.ReplaceAll(pattern, "+", " + ")
	pattern = strings.ReplaceAll(pattern, "-", " - ")
	pattern = strings.ReplaceAll(pattern, "*", " * ")
	pattern = strings.ReplaceAll(pattern, "/", " / ")
	pattern = strings.ReplaceAll(pattern, "(", " ( ")
	pattern = strings.ReplaceAll(pattern, ")", " ) ")
	lst = strings.Fields(pattern)
	for _, v := range lst {
		if unicode.IsDigit([]rune(v)[0]) {
			ouput_stack = append(ouput_stack, v)
		} else if precedence, exists := operators[v]; exists {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if topPrecedence, topExists := operators[top]; topExists && topPrecedence >= precedence {
					ouput_stack = append(ouput_stack, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, v)
		} else if v == "(" {
			stack = append(stack, v)
		} else if v == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				ouput_stack = append(ouput_stack, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("too many closing parentheses")
			}
			stack = stack[:len(stack)-1]
		}
	}
	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("too many opening parentheses")
		}
		ouput_stack = append(ouput_stack, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return ouput_stack, nil
}
