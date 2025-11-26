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

func (st *ServerTmpl) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input store.Event
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest) // 400
		log.Println("Method:", r.Method)
		log.Println("Content-Type:", r.Header.Get("Content-Type"))
	}

	createdEvent, err := st.Service.CreateEvent(&input)
	if err != nil {
		log.Printf("ошибка при создании события %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable) // 503 хотя мне кажется 404

	}

	fmt.Printf("EventId : %d | User : %d | Date : %s | Name : %s\n", createdEvent.ID, createdEvent.UserID, createdEvent.Date, createdEvent.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201

	json.NewEncoder(w).Encode(createdEvent)
}

func (st *ServerTmpl) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest) // 400
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest) // 400
		return
	}

	var input store.Event
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest) // 400
		log.Println("Method:", r.Method)
		log.Println("Content-Type:", r.Header.Get("Content-Type"))
	}

	updatedEvent, err := st.Service.UpdateEvent(id, &input)
	if err != nil {
		log.Printf("ошибка при обновления события %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable) // 503 хотя мне кажется 404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200

	json.NewEncoder(w).Encode(updatedEvent)
}

func (st *ServerTmpl) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest) // 400
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest) // 400
		return
	}

	deletedEv, err := st.Service.DeleteEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable) // 503 если не нашли событие
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200

	json.NewEncoder(w).Encode(deletedEv)
}

func (st *ServerTmpl) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userIDS := query.Get("user_id")
	dateStr := query.Get("date")

	if userIDS == " " || dateStr == " " {
		http.Error(w, "missing parametrs", http.StatusBadRequest) // 400
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest) // 400
		return
	}

	eventsForDay, err := st.Service.GetEventsForDay(userID, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
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
		http.Error(w, "missing parameters", http.StatusBadRequest) // 400
		return
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest) // 400
		return
	}

	eventsForWeek, err := st.Service.GetEventsForWeek(userID, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200

	json.NewEncoder(w).Encode(eventsForWeek)
}

func (st *ServerTmpl) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userIDS := query.Get("user_id")
	dateStr := query.Get("date")

	if userIDS == "" || dateStr == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest) // 400
		return
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest) // 400
		return
	}

	eventsForMonth, err := st.Service.GetEventsForMonth(userID, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200

	json.NewEncoder(w).Encode(eventsForMonth)
}
