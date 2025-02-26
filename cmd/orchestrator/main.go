package main

import (
	"log"
	"net/http"

	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	api := orchestrator.NewAPI()
	r.HandleFunc("/api/v1/calculate", api.AddNewExpression)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
