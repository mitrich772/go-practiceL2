package store

import (
	"fmt"
	"sync"
)

// Event описывает одно событие
// @Description Событие
type Event struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Date   string `json:"date"`
	Name   string `json:"name"`
}

// Store интерфейс для работы с событиями
type Store interface {
	Create(*Event)                        // добавление нового события
	Update(int64, *Event) (*Event, error) // обновление события по ID, ошибка если не найдено
	Delete(int64) (*Event, error)         // удаление события по ID, возвращает удалённое
	Get(int64) (*Event, error)            // получение события по ID
	List(int64) []*Event                  // список событий конкретного пользователя
}

// EventStore потокобезопасное хранилище событий
type EventStore struct {
	data map[int64]*Event
	mu   sync.Mutex
}

// NewEventStore создает новый EventStore
func NewEventStore() *EventStore {
	return &EventStore{data: make(map[int64]*Event)}
}

// Create добавляет событие в хранилище
func (es *EventStore) Create(e *Event) {
	es.mu.Lock()
	es.data[e.ID] = e
	es.mu.Unlock()
}

// Update обновляет существующий event по id
// Если события нет то возвращает ошибку
func (es *EventStore) Update(id int64, eNew *Event) (*Event, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	if _, exists := es.data[id]; !exists {
		return nil, fmt.Errorf("event with id=%d not found", id)
	}

	eNew.ID = id // сохраняем ID
	es.data[id] = eNew

	return eNew, nil
}

// Delete удаляет событие по ID, возвращает удалённое событие
// Если события нет то возвращает ошибку
func (es *EventStore) Delete(id int64) (*Event, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	ev, exists := es.data[id]
	if !exists {
		return nil, fmt.Errorf("event with id=%d not found", id)
	}

	delete(es.data, id)
	return ev, nil
}

// Get возвращает событие по ID, если нет - ошибка
func (es *EventStore) Get(id int64) (*Event, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	ev, exists := es.data[id]
	if !exists {
		return nil, fmt.Errorf("event with id=%d not found", id)
	}

	return ev, nil
}

// List возвращает все события для переданного пользователя (userID)
func (es *EventStore) List(userID int64) []*Event {
	es.mu.Lock()
	defer es.mu.Unlock()

	res := make([]*Event, 0)
	for _, ev := range es.data {
		if ev.UserID == userID {
			res = append(res, ev)
		}
	}
	return res
}
