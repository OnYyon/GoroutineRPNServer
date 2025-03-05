package main

import (
	"github.com/OnYyon/GoroutineRPNServer/cmd/agent"
	"github.com/OnYyon/GoroutineRPNServer/cmd/orchestrator"
	"github.com/OnYyon/GoroutineRPNServer/web"
)

func main() {
	agent.StartAgent()
	orchestrator.StartOrchestrator()
	web.StartWeb()
}
