package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

// Структура для выражения
type Expression struct {
	ID     string  `json:"id"`
	Input  string  `json:"input"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

// Главная страница
func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

// Отправка выражения на сервер
func addExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	expr := r.FormValue("expression")
	fmt.Println(expr)
	data := map[string]string{"expression": expr}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://127.0.0.1:8080/api/v1/calculate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Получение списка выражений
func getExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://127.0.0.1:8080/api/v1/expressions")
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             // Замените на ваш источник
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")    // Укажите методы, которые вы хотите разрешить
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Укажите заголовки, которые вы хотите разрешить
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add", addExpressionHandler)
	http.HandleFunc("/expressions", getExpressionsHandler)

	// Подключение статических файлов (CSS, JS)
	http.Handle("/static/", corsMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))))

	http.ListenAndServe(":8081", nil)
}
