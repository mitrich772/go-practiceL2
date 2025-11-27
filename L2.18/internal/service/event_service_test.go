package service

import (
	"testing"

	"github.com/golang/mock/gomock"

	"L2.18/internal/store"
	"L2.18/internal/store/mocks"
)

func TestEventService_GetEventsForDay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	srv := NewEventService(mockStore)

	userID := int64(1)
	mockEvents := []*store.Event{
		{ID: 1, UserID: userID, Date: "2025-11-27"},
		{ID: 2, UserID: userID, Date: "2025-11-28"},
	}

	mockStore.EXPECT().List(userID).Return(mockEvents)

	events, err := srv.GetEventsForDay(userID, "2025-11-27")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 1 || events[0].ID != 1 {
		t.Errorf("expected single event with ID==1, got %+v", events)
	}
}

func TestEventService_GetEventsForWeek(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	srv := NewEventService(mockStore)

	userID := int64(1)
	mockEvents := []*store.Event{
		{ID: 1, UserID: userID, Date: "2025-11-25"}, // week 48
		{ID: 2, UserID: userID, Date: "2025-11-27"}, // week 48
		{ID: 3, UserID: userID, Date: "2025-12-01"}, // week 49
	}

	mockStore.EXPECT().List(userID).Return(mockEvents)

	events, err := srv.GetEventsForWeek(userID, "2025-11-27") // week 48
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("expected 2 events for week, got %d", len(events))
	}
}

func TestEventService_GetEventsForMonth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	srv := NewEventService(mockStore)

	userID := int64(1)
	mockEvents := []*store.Event{
		{ID: 1, UserID: userID, Date: "2025-11-25"},
		{ID: 2, UserID: userID, Date: "2025-11-27"},
		{ID: 3, UserID: userID, Date: "2025-12-01"},
	}

	mockStore.EXPECT().List(userID).Return(mockEvents)

	events, err := srv.GetEventsForMonth(userID, "2025-11-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("expected 2 events for month, got %d", len(events))
	}
}

func TestEventService_CreateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	srv := NewEventService(mockStore)

	// создаём событие без ID, ID должен присвоиться в CreateEvent
	inputEvent := &store.Event{
		UserID: 1,
		Date:   "2025-11-27",
		Name:   "Test Event",
	}

	// ожидаем, что Store.Create будет вызвано с событием, у которого ID == 1 (service сам присвоит)
	mockStore.EXPECT().Create(gomock.AssignableToTypeOf(&store.Event{})).DoAndReturn(
		func(e *store.Event) {
			if e.ID != 1 {
				t.Errorf("expected ID=1, got %d", e.ID)
			}
			if e.UserID != inputEvent.UserID || e.Date != inputEvent.Date || e.Name != inputEvent.Name {
				t.Errorf("event fields mismatch: got %+v", e)
			}
		},
	)

	created, err := srv.CreateEvent(inputEvent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("expected ID=1, got %d", created.ID)
	}

	if created.UserID != inputEvent.UserID || created.Date != inputEvent.Date || created.Name != inputEvent.Name {
		t.Errorf("created event fields mismatch: %+v", created)
	}
}
