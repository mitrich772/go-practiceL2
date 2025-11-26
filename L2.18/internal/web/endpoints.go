package web

import (
	"strconv"

	"L2.18/internal/service"
	"L2.18/internal/store"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ServerTmpl struct {
	Service *service.EventService
}

func (st *ServerTmpl) CreateEvent(w http.ResponseWriter, r *http.Request) { // POST json
	var input store.Event
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		log.Println("Method:", r.Method)
		log.Println("Content-Type:", r.Header.Get("Content-Type"))
	}

	createdEvent, err := st.Service.CreateEvent(&input)
	if err != nil {
		log.Printf("ошибка при создании события %v", err)
	}

	fmt.Printf("EventId : %d | User : %d | Date : %s | Name : %s\n", createdEvent.ID, createdEvent.UserID, createdEvent.Date, createdEvent.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201

	json.NewEncoder(w).Encode(createdEvent)
}

func (st *ServerTmpl) DeleteEvent(w http.ResponseWriter, r *http.Request) { // DELETE /events_for_day/123
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = st.Service.DeleteEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 без тела
}

func (st *ServerTmpl) GetEventsForDay(w http.ResponseWriter, r *http.Request) { // GET /events_for_day?user_id=1&date=2023-12-31
	query := r.URL.Query()

	userIDS := query.Get("user_id")
	dateStr := query.Get("date")

	if userIDS == " " || dateStr == " " {
		http.Error(w, "missing parametrs", http.StatusBadRequest)
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	eventsForDay, err := st.Service.GetEventsForDay(userID, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200

	json.NewEncoder(w).Encode(eventsForDay)
}

func (st *ServerTmpl) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userIDS := query.Get("user_id")
	dateStr := query.Get("date")

	if userIDS == "" || dateStr == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	eventsForWeek, err := st.Service.GetEventsForWeek(userID, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(eventsForWeek)
}

func (st *ServerTmpl) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userIDS := query.Get("user_id")
	dateStr := query.Get("date")

	if userIDS == "" || dateStr == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	eventsForMonth, err := st.Service.GetEventsForMonth(userID, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(eventsForMonth)
}
