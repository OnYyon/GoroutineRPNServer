package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/OnYyon/GoroutineRPNServer/iternal/agent"
	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := config.NewConfig()
	fmt.Println(cfg)
	var wg sync.WaitGroup
	for i := 0; i < cfg.ComputerPower; i++ {
		wg.Add(1)
		fmt.Printf("Starting gouroutine: %v\n", i)
		go agent.Worker(&wg, i)
	}

	wg.Wait()
}
