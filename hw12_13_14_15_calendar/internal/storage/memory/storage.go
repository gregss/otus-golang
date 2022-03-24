package memorystorage

import (
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// mu     sync.RWMutex
	events []storage.Event
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Append(event storage.Event) (int, error) {
	s.events = append(s.events, event)
	return event.ID, nil
}

func (s *Storage) Edit(id int, newEvent storage.Event) error {
	for i, event := range s.events {
		if event.ID == id {
			s.events[i] = newEvent
			break
		}
	}

	return nil
}

func (s *Storage) Delete(id int) error {
	for i, event := range s.events {
		if event.ID == id {
			s.events = append(s.events[:i], s.events[i+1:]...)
			break
		}
	}

	return nil
}

func (s *Storage) DayEvents(date time.Time) []storage.Event {
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.Time.Day() == date.Day() {
			result = append(result, event)
		}
	}

	return result
}

func (s *Storage) WeekEvents(date time.Time) []storage.Event {
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.Time.Weekday() == date.Weekday() {
			result = append(result, event)
		}
	}

	return result
}

func (s *Storage) MonthEvents(date time.Time) []storage.Event {
	result := make([]storage.Event, 0)
	for _, event := range s.events {
		if event.Time.Month() == date.Month() {
			result = append(result, event)
		}
	}

	return result
}
