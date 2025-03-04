package main

import (
	"log"
	"net/http"

	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	r := mux.NewRouter()

	api := orchestrator.NewAPI()
	r.HandleFunc("/api/v1/calculate", api.AddNewExpression).Methods("POST")
	r.HandleFunc("/api/v1/expressions", api.GetSliceOfExpressions).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", api.GetExpressionByID).Methods("GET")

	//Iternal handlers
	r.HandleFunc("/iternal/task", api.GetTasksToAgent).Methods("GET")
	r.HandleFunc("/iternal/task", api.GetPostResult).Methods("POST")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
