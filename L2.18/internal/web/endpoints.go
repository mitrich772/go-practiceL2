package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"L2.18/internal/service"
	"L2.18/internal/store"
)

// EventServer реализует HTTP-обработчики для работы с событиями
type EventServer struct {
	Service *service.EventService
}

// CreateEvent создаёт новое событие
// @Summary      Create event
// @Description  Creates a new event
// @Accept       json
// @Produce      json
// @Success      201 {object} store.Event
// @Router       /create_event [post]
func (st *EventServer) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input store.Event
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest) // 400
		log.Println("Method:", r.Method)
		log.Println("Content-Type:", r.Header.Get("Content-Type"))
	}

	if input.Name == "" {
		http.Error(w, "name is empty", http.StatusBadRequest)
		return
	}

	if input.Date == "" {
		http.Error(w, "date is empty", http.StatusBadRequest)
		return
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

// UpdateEvent обновляет существующее событие по ID
// @Summary      Update event
// @Description  Updates an existing event by id
// @Accept       json
// @Produce      json
// @Success      200 {object} store.Event
// @Router       /update_event [post]
func (st *EventServer) UpdateEvent(w http.ResponseWriter, r *http.Request) {
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

// DeleteEvent удаляет событие по ID
// @Summary      Delete event
// @Description  Deletes an event by id
// @Produce      json
// @Success      200 {object} store.Event
// @Router       /delete_event [delete]
func (st *EventServer) DeleteEvent(w http.ResponseWriter, r *http.Request) {
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

// GetEventsForDay возвращает события пользователя за конкретный день
// @Summary      Get events for day
// @Description  Returns events of user for specific day
// @Produce      json
// @Success      200 {array} store.Event
// @Router       /events_for_day [get]
func (st *EventServer) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
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

// GetEventsForWeek возвращает события пользователя за неделю, содержащую указанную дату
// @Summary      Get events for week
// @Description  Returns events of user for week of given date
// @Produce      json
// @Success      200 {array} store.Event
// @Router       /events_for_week [get]
func (st *EventServer) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
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

// GetEventsForMonth возвращает события пользователя за месяц, содержащий указанную дату
// @Summary      Get events for month
// @Description  Returns events of user for month of given date
// @Produce      json
// @Success      200 {array} store.Event
// @Router       /events_for_month [get]
func (st *EventServer) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
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
