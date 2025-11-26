package store

import (
	"fmt"
	"sync"
)

type Event struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Date   string `json:"date"`
	Name   string `json:"name"`
}

type Store interface {
	Create(Event)
	Update(Event) (Event, error)
	Delete(int64) error
	Get(int64) (Event, error)
	List(int64) []Event
}

type EventStore struct {
	data map[int64]Event
	mu   sync.Mutex
}

func NewEventStore() *EventStore {
	return &EventStore{data: make(map[int64]Event)}
}
func (es *EventStore) Create(e Event) {
	es.mu.Lock()
	es.data[e.ID] = e
	es.mu.Unlock()
}

func (es *EventStore) Update(e Event) (Event, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	if _, exists := es.data[e.ID]; !exists {
		return Event{}, fmt.Errorf("no event with id=%d to uptade", e.ID)
	}

	es.data[e.ID] = e
	return e, nil
}

func (es *EventStore) Delete(id int64) error {
	es.mu.Lock()
	defer es.mu.Unlock()

	if _, exists := es.data[id]; !exists {
		return fmt.Errorf("no id=%d to delete", id)
	}

	delete(es.data, id)
	return nil
}

func (es *EventStore) Get(id int64) (Event, error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	if _, exists := es.data[id]; !exists {
		return Event{}, fmt.Errorf("no id=%d to delete", id)
	}
	return es.data[id], nil
}

// List возращает все события для переданного userId
func (es *EventStore) List(userId int64) []Event {
	es.mu.Lock()
	defer es.mu.Unlock()

	res := make([]Event, 0, len(es.data))
	for _, ev := range es.data {
		if ev.UserID == userId {
			res = append(res, ev)
		}
	}
	return res
}
