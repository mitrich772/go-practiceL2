package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "L2.18/docs"
	"L2.18/internal/logging"
	"L2.18/internal/service"
	"L2.18/internal/store"
	"L2.18/internal/web"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	var port string
	flag.StringVar(&port, "p", "1235", "set port for server")
	flag.Parse()

	muxServ := http.NewServeMux()

	evServ := service.NewEventService(store.NewEventStore())

	stm := web.EventServer{Service: evServ}

	muxServ.HandleFunc("/create_event", logging.Middleware(stm.CreateEvent))
	muxServ.HandleFunc("/update_event", logging.Middleware(stm.UpdateEvent))
	muxServ.HandleFunc("/delete_event", logging.Middleware(stm.DeleteEvent))
	muxServ.HandleFunc("/events_for_day", logging.Middleware(stm.GetEventsForDay))
	muxServ.HandleFunc("/events_for_week", logging.Middleware(stm.GetEventsForWeek))
	muxServ.HandleFunc("/events_for_month", logging.Middleware(stm.GetEventsForMonth))
	// Swagger UI
	muxServ.Handle("/swagger/", httpSwagger.WrapHandler)

	go func() {
		log.Println("Web сервер запущен на порту", port)
		if err := http.ListenAndServe(":"+port, muxServ); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
