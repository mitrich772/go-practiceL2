package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"L2.18/internal/service"
	"L2.18/internal/store"
	"L2.18/internal/web"
)

func main() {
	muxServ := http.NewServeMux()

	evServ := service.NewEventService(store.NewEventStore())

	stm := web.ServerTmpl{Service: evServ}

	muxServ.HandleFunc("/create_event", stm.CreateEvent)
	muxServ.HandleFunc("/events_for_day", stm.GetEventsForDay)
	muxServ.HandleFunc("/events_for_week", stm.GetEventsForWeek)
	muxServ.HandleFunc("/events_for_mounth", stm.GetEventsForMonth)
	port := "1235"
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
