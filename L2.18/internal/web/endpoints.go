package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"L2.18/internal/service"
	"L2.18/internal/store"
)

type EventServer struct {
	Service *service.EventService
}

// CreateEvent
// @Summary      Создать событие
// @Description  Создаёт новое событие
// @Accept       json
// @Produce      json
// @Param        input body store.Event true "JSON события"
// @Success      201 {object} web.ResultResponse "Созданное событие в поле result"
// @Failure      400 {object} web.ErrorResponse "Некорректные данные"
// @Failure      503 {object} web.ErrorResponse "Ошибка логики"
// @Router       /create_event [post]
func (st *EventServer) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input store.Event
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "невалидный JSON")
		return
	}

	if input.Name == "" {
		writeError(w, http.StatusBadRequest, "имя не может быть пустым")
		return
	}

	if input.Date == "" {
		writeError(w, http.StatusBadRequest, "дата не может быть пустой")
		return
	}

	createdEvent, err := st.Service.CreateEvent(&input)
	if err != nil {
		log.Printf("ошибка при создании события: %v", err)
		writeError(w, http.StatusServiceUnavailable, "ошибка сервера")
		return
	}

	writeResult(w, http.StatusCreated, createdEvent)
}

// UpdateEvent
// @Summary      Обновить событие
// @Description  Обновляет событие по указанному ID
// @Accept       json
// @Produce      json
// @Param        id query int true "ID события"
// @Param        input body store.Event true "JSON события"
// @Success      200 {object} web.ResultResponse "Обновлённое событие"
// @Failure      400 {object} web.ErrorResponse "Некорректные параметры"
// @Failure      404 {object} web.ErrorResponse "Событие не найдено"
// @Failure      503 {object} web.ErrorResponse "Ошибка логики"
// @Router       /update_event [post]
func (st *EventServer) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "не указан id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неправильный формат id")
		return
	}

	var input store.Event
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "невалидный JSON")
		return
	}

	updatedEvent, err := st.Service.UpdateEvent(id, &input)
	if err != nil {
		log.Printf("ошибка при обновлении: %v", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeResult(w, http.StatusOK, updatedEvent)
}

// DeleteEvent
// @Summary      Удалить событие
// @Description  Удаляет событие по ID
// @Produce      json
// @Param        id query int true "ID события"
// @Success      200 {object} web.ResultResponse "Удалённое событие"
// @Failure      400 {object} web.ErrorResponse "Некорректные параметры"
// @Failure      404 {object} web.ErrorResponse "Событие не найдено"
// @Failure      503 {object} web.ErrorResponse "Ошибка логики"
// @Router       /delete_event [delete]
func (st *EventServer) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "не указан id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неправильный формат id")
		return
	}

	deletedEv, err := st.Service.DeleteEvent(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeResult(w, http.StatusOK, deletedEv)
}

// GetEventsForDay
// @Summary      Получить события за день
// @Description  Возвращает события пользователя за указанный день
// @Produce      json
// @Param        user_id query int true "ID пользователя"
// @Param        date query string true "Дата YYYY-MM-DD"
// @Success      200 {object} web.ResultResponse "События за день"
// @Failure      400 {object} web.ErrorResponse "Некорректные параметры"
// @Failure      503 {object} web.ErrorResponse "Ошибка логики"
// @Router       /events_for_day [get]
func (st *EventServer) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userIDS := q.Get("user_id")
	dateStr := q.Get("date")

	if userIDS == "" || dateStr == "" {
		writeError(w, http.StatusBadRequest, "не указаны параметры")
		return
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неправильный формат user_id")
		return
	}

	events, err := st.Service.GetEventsForDay(userID, dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeResult(w, http.StatusOK, events)
}

// GetEventsForWeek
// @Summary      Получить события за неделю
// @Description  Возвращает события пользователя за неделю
// @Produce      json
// @Param        user_id query int true "ID пользователя"
// @Param        date query string true "Дата YYYY-MM-DD"
// @Success      200 {object} web.ResultResponse "События за неделю"
// @Failure      400 {object} web.ErrorResponse "Некорректные параметры"
// @Failure      503 {object} web.ErrorResponse "Ошибка логики"
// @Router       /events_for_week [get]
func (st *EventServer) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userIDS := q.Get("user_id")
	dateStr := q.Get("date")

	if userIDS == "" || dateStr == "" {
		writeError(w, http.StatusBadRequest, "не указаны параметры")
		return
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неправильный формат user_id")
		return
	}

	events, err := st.Service.GetEventsForWeek(userID, dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeResult(w, http.StatusOK, events)
}

// GetEventsForMonth
// @Summary      Получить события за месяц
// @Description  Возвращает события пользователя за месяц
// @Produce      json
// @Param        user_id query int true "ID пользователя"
// @Param        date query string true "Дата YYYY-MM-DD"
// @Success      200 {object} web.ResultResponse "События за месяц"
// @Failure      400 {object} web.ErrorResponse "Некорректные параметры"
// @Failure      503 {object} web.ErrorResponse "Ошибка логики"
// @Router       /events_for_month [get]
func (st *EventServer) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userIDS := q.Get("user_id")
	dateStr := q.Get("date")

	if userIDS == "" || dateStr == "" {
		writeError(w, http.StatusBadRequest, "не указаны параметры")
		return
	}

	userID, err := strconv.ParseInt(userIDS, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неправильный формат user_id")
		return
	}

	events, err := st.Service.GetEventsForMonth(userID, dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeResult(w, http.StatusOK, events)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

func writeResult(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ResultResponse{Result: data})
}
