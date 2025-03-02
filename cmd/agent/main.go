package main

import (
	"github.com/OnYyon/GoroutineRPNServer/iternal/agent"
	"github.com/OnYyon/GoroutineRPNServer/iternal/config"
	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
)

func main() {
	cfg := config.NewConfig()
	computingPower := cfg.ComputerPower

	for i := 0; i < computingPower; i++ {
		go agent.Worker(&orchestrator.Wg)
	}
}
