package memorystorage_test

import (
	"testing"
	"time"

	storage "github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	s := memorystorage.New()

	now := time.Date(2022, 0o3, 0o1, 10, 0, 0, 0, time.UTC)

	id, _ := s.Append(storage.Event{
		ID:          1,
		Title:       "title1",
		Time:        now,
		Duration:    100 * time.Second,
		Description: "description1",
		UserID:      1,
		NotifyTime:  now,
	})

	require.Equal(t, id, 1)

	_ = s.Edit(1, storage.Event{
		ID:          2,
		Title:       "title2",
		Time:        now,
		Duration:    200 * time.Second,
		Description: "description2",
		UserID:      2,
		NotifyTime:  now,
	})

	events := s.DayEvents(now)
	require.Equal(t, 1, len(events))
	require.Equal(t, 2, events[0].ID)

	events = s.DayEvents(time.Now())
	require.Equal(t, 0, len(events))

	_ = s.Delete(2)
	events = s.DayEvents(now)
	require.Equal(t, 0, len(events))
}
