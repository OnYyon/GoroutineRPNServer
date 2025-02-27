package main

import (
	"fmt"
	"time"
)

const (
	TIME_ADDITION_MS        = 100 // Пример значения, замените на реальное измерение
	TIME_SUBTRACTION_MS     = 100 // Пример значения, замените на реальное измерение
	TIME_MULTIPLICATIONS_MS = 100 // Пример значения, замените на реальное измерение
	TIME_DIVISIONS_MS       = 100 // Пример значения, замените на реальное измерение
)

func measureTimeAddition() int64 {
	start := time.Now()
	// Пример операции сложения
	a := 1000000
	b := 2000000
	for i := 0; i < 1000000; i++ {
		_ = a + b
	}
	duration := time.Since(start)
	return duration.Milliseconds()
}

func measureTimeSubtraction() int64 {
	start := time.Now()
	// Пример операции вычитания
	a := 1000000
	b := 2000000
	for i := 0; i < 1000000; i++ {
		_ = a - b
	}
	duration := time.Since(start)
	return duration.Milliseconds()
}

func measureTimeMultiplication() int64 {
	start := time.Now()
	// Пример операции умножения
	a := 1000000
	b := 2000000
	for i := 0; i < 1000000; i++ {
		_ = a * b
	}
	duration := time.Since(start)
	return duration.Milliseconds()
}

func measureTimeDivision() int64 {
	start := time.Now()
	// Пример операции деления
	a := 1000000
	b := 2000000
	for i := 0; i < 1000000; i++ {
		_ = a / b
	}
	duration := time.Since(start)
	return duration.Milliseconds()
}

func main() {
	additionTime := measureTimeAddition()
	subtractionTime := measureTimeSubtraction()
	multiplicationTime := measureTimeMultiplication()
	divisionTime := measureTimeDivision()

	fmt.Printf("Time for addition: %d ms\n", additionTime)
	fmt.Printf("Time for subtraction: %d ms\n", subtractionTime)
	fmt.Printf("Time for multiplication: %d ms\n", multiplicationTime)
	fmt.Printf("Time for division: %d ms\n", divisionTime)
}
