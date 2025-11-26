package service

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"L2.18/internal/store"
)

type EventService struct {
	Storage store.Store
	lastID  int64
}

func NewEventService(store store.Store) *EventService {
	return &EventService{
		Storage: store,
		lastID:  0,
	}
}

func (srv *EventService) CreateEvent(event *store.Event) (*store.Event, error) {
	event.ID = atomic.AddInt64(&srv.lastID, 1)
	srv.Storage.Create(*event)
	return event, nil
}

func (srv *EventService) UpdateEvent(event *store.Event) (*store.Event, error) {
	updated, err := srv.Storage.Update(*event)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (srv *EventService) DeleteEvent(id int64) error {
	return srv.Storage.Delete(id)
}

func (srv *EventService) GetEventsForDay(userId int64, dayStr string) ([]store.Event, error) {
	events := srv.Storage.List(userId)

	var res []store.Event
	const layout = "2006-01-02" // формат даты в строке

	day, err := time.Parse(layout, dayStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты для поиска: %s", dayStr)
	}

	for _, ev := range events {
		if ev.UserID != userId {
			continue
		}

		evDate, err := time.Parse(layout, ev.Date)
		if err != nil {
			log.Printf("не удалось распарсить дату %s | %v", evDate, err)
			continue // пропускаем некорректные даты
		}

		if evDate.Year() == day.Year() &&
			evDate.Month() == day.Month() &&
			evDate.Day() == day.Day() {
			res = append(res, ev)
		}
	}

	return res, nil
}
func (srv *EventService) GetEventsForWeek(userId int64, dayStr string) ([]store.Event, error) {
	events := srv.Storage.List(userId)

	var res []store.Event
	const layout = "2006-01-02"

	day, err := time.Parse(layout, dayStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты для поиска: %s", dayStr)
	}

	year, week := day.ISOWeek()

	for _, ev := range events {
		if ev.UserID != userId {
			continue
		}

		evDate, err := time.Parse(layout, ev.Date)
		if err != nil {
			log.Printf("не удалось распарсить дату %s", evDate)
			continue // пропускаем некорректные даты
		}

		evYear, evWeek := evDate.ISOWeek()

		if evYear == year && evWeek == week {
			res = append(res, ev)
		}
	}

	return res, nil
}

func (srv *EventService) GetEventsForMonth(userId int64, dayStr string) ([]store.Event, error) {
	events := srv.Storage.List(userId)

	var res []store.Event
	const layout = "2006-01-02"

	day, err := time.Parse(layout, dayStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты для поиска: %s", dayStr)
	}

	targetYear := day.Year()
	targetMonth := day.Month()

	for _, ev := range events {
		if ev.UserID != userId {
			continue
		}

		evDate, err := time.Parse(layout, ev.Date)
		if err != nil {
			log.Printf("не удалось распарсить дату %s", evDate)
			continue
		}

		if evDate.Year() == targetYear && evDate.Month() == targetMonth {
			res = append(res, ev)
		}
	}

	return res, nil
}
