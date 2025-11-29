package service

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"L2.18/internal/store"
)

// EventService предоставляет методы работы с Store
type EventService struct {
	Storage store.Store
	lastID  int64
}

// NewEventService создаёт новый сервис событий с заданным хранилищем
func NewEventService(s store.Store) *EventService {
	return &EventService{
		Storage: s,
		lastID:  0,
	}
}

// CreateEvent добавляет новое событие и присваивает ему уникальный ID
func (srv *EventService) CreateEvent(event *store.Event) (*store.Event, error) {
	// проверяем корректность даты
	_, err := time.Parse("2006-01-02", event.Date)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты, ожидается YYYY-MM-DD: %w", err)
	}

	event.ID = atomic.AddInt64(&srv.lastID, 1)
	srv.Storage.Create(event)
	return event, nil
}

// UpdateEvent обновляет существующее событие по ID на переданное в event
func (srv *EventService) UpdateEvent(id int64, event *store.Event) (*store.Event, error) {
	return srv.Storage.Update(id, event)
}

// DeleteEvent удаляет событие по ID
func (srv *EventService) DeleteEvent(id int64) (*store.Event, error) {
	return srv.Storage.Delete(id)
}

// GetEventsForDay возвращает события пользователя за конкретный день
func (srv *EventService) GetEventsForDay(userID int64, dayStr string) ([]*store.Event, error) {
	events := srv.Storage.List(userID)

	const layout = "2006-01-02"
	day, err := time.Parse(layout, dayStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты: %s", dayStr)
	}

	var res []*store.Event
	for _, ev := range events {
		evDate, err := time.Parse(layout, ev.Date)
		if err != nil {
			log.Printf("не удалось распарсить дату %s | %v", ev.Date, err)
			continue
		}

		if evDate.Year() == day.Year() &&
			evDate.Month() == day.Month() &&
			evDate.Day() == day.Day() {
			res = append(res, ev)
		}
	}

	return res, nil
}

// GetEventsForWeek возвращает события пользователя за неделю, содержащую указанную дату
func (srv *EventService) GetEventsForWeek(userID int64, dayStr string) ([]*store.Event, error) {
	events := srv.Storage.List(userID)

	const layout = "2006-01-02"
	day, err := time.Parse(layout, dayStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты: %s", dayStr)
	}

	year, week := day.ISOWeek()
	var res []*store.Event
	for _, ev := range events {
		evDate, err := time.Parse(layout, ev.Date)
		if err != nil {
			log.Printf("не удалось распарсить дату %s", ev.Date)
			continue
		}

		evYear, evWeek := evDate.ISOWeek()
		if evYear == year && evWeek == week {
			res = append(res, ev)
		}
	}

	return res, nil
}

// GetEventsForMonth возвращает события пользователя за месяц, содержащий указанную дату
func (srv *EventService) GetEventsForMonth(userID int64, dayStr string) ([]*store.Event, error) {
	events := srv.Storage.List(userID)

	const layout = "2006-01-02"
	day, err := time.Parse(layout, dayStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты: %s", dayStr)
	}

	targetYear := day.Year()
	targetMonth := day.Month()

	var res []*store.Event
	for _, ev := range events {
		evDate, err := time.Parse(layout, ev.Date)
		if err != nil {
			log.Printf("не удалось распарсить дату %s", ev.Date)
			continue
		}

		if evDate.Year() == targetYear && evDate.Month() == targetMonth {
			res = append(res, ev)
		}
	}

	return res, nil
}
