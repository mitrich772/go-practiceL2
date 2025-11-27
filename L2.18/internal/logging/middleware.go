package logging

import (
	"log"
	"net/http"
	"time"
)

// Middleware оборачивает HTTP-обработчик и логирует:
// - метод и путь запроса
// - время выполнения обработчика
func Middleware(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("%s %s", r.Method, r.URL.Path)

		handler(w, r)

		log.Printf("done in %v\n", time.Since(start))
	}
}
