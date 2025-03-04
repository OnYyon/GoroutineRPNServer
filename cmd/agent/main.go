package main

import (
	"fmt"
	"sync"

	"github.com/OnYyon/GoroutineRPNServer/iternal/agent"
	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
)

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
