package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TimeAdd         time.Duration
	TimeSubtraction time.Duration
	TimeMultiply    time.Duration
	TimeDivision    time.Duration
	ComputerPower   int
}

func NewConfig() *Config {
	return &Config{
		TimeAdd:         getEnvTime("TIME_ADDITION_MS", 9999*time.Millisecond),
		TimeSubtraction: getEnvTime("TIME_SUBTRACTION_MS", 9999*time.Millisecond),
		TimeMultiply:    getEnvTime("TIME_MULTIPLICATIONS_MS", 9999*time.Millisecond),
		TimeDivision:    getEnvTime("TIME_DIVISIONS_MS", 9999*time.Millisecond),
		ComputerPower:   getEnvInt("COMPUTING_POWER", 3),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	fmt.Println(valueStr)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvTime(name string, defaultValueMs time.Duration) time.Duration {
	valueStr := getEnv(name, "")
	if timeInt, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return time.Duration(timeInt) * time.Millisecond
	}
	return defaultValueMs
}
