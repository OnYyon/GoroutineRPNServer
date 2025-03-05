package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OnYyon/GoroutineRPNServer/iternal/orchestrator"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
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

	fmt.Println("Starting orhcestrator on localhost:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}
