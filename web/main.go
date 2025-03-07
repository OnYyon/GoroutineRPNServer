package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

var templates = template.Must(template.ParseGlob("./web/templates/*.html"))

type Expression struct {
	ID     string  `json:"id"`
	Input  string  `json:"input"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func addExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	expr := r.FormValue("expression")
	fmt.Println(expr)
	data := map[string]string{"expression": expr}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://orchestrator:8080/api/v1/calculate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://orchestrator:8080/api/v1/expressions")
	if err != nil {
		http.Error(w, "Failed to fetch expressions", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add", addExpressionHandler)
	http.HandleFunc("/expressions", getExpressionsHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	fmt.Println("Starting web-interface on localhost:8081")
	http.ListenAndServe("0.0.0.0:8081", nil)
}
