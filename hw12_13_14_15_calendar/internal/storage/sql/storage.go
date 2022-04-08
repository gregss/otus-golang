package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gregss/otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type DBStorage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type Storage struct {
	dsn string
	con *sqlx.DB
}

func New(dsn string) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.con, err = sqlx.Open("pgx", s.dsn)

	if err != nil {
		log.Fatalf("err %v", err)
	}

	if s.con == nil {
		log.Fatal("nil connection")
	}

	return nil
}

func (s *Storage) Ping() error {
	return s.con.Ping()
}

func (s *Storage) Close(ctx context.Context) error {
	return s.con.Close()
}

func (s *Storage) Append(event storage.Event) (int, error) {
	query := `insert into events(owner, title, descr, start_date, end_date) values($1, $2, $3, $4, $5)`

	result, _ := s.con.Exec(query,
		event.UserID,
		event.Title,
		event.Description,
		event.Time,
		event.Time.Add(event.Duration),
	)

	eventID, _ := result.LastInsertId()
	return int(eventID), nil
}

func (s *Storage) Edit(id int, newEvent storage.Event) error {
	query := `UPDATE events SET owner = $1, title = $2, descr = $3, start_date = $4, end_date = $5 where id = $6`
	_, err := s.con.Exec(query,
		newEvent.UserID,
		newEvent.Title,
		newEvent.Description,
		newEvent.Time,
		newEvent.Time.Add(newEvent.Duration),
		id,
	)
	return err
}

func (s *Storage) Delete(id int) error {
	_, err := s.con.Exec(`delete from events where id = $1`, id)
	return err
}

func (s *Storage) DayEvents(date time.Time) []storage.Event {
	events := make([]storage.Event, 0)
	rows, _ := s.con.Query(`select * from events where start_date = $1`, date)
	defer rows.Close()
	for rows.Next() {
		event := new(storage.Event)

		var startDate, endDate time.Time
		var startTime, endTime sql.NullString
		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Title,
			&event.Description,
			&startDate,
			&startTime,
			&endDate,
			&endTime,
		)
		if err != nil {
			fmt.Printf("%v", err)
		}

		events = append(events, *event)
	}

	if rows.Err() != nil {
		fmt.Printf("%v", rows.Err())
		return nil
	}

	return events
}

func (s *Storage) WeekEvents(date time.Time) []storage.Event {
	return nil // todo
}

func (s *Storage) MonthEvents(date time.Time) []storage.Event {
	return nil // todo
}
