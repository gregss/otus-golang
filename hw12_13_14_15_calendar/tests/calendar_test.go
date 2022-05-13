//go:build integrations
// +build integrations

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	internalhttp "github.com/gregss/otus/hw12_13_14_15_calendar/internal/server/http"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/suite"
)

type CalendarSuite struct {
	suite.Suite
}

func (s *CalendarSuite) SetupSuite() {
}

func (s *CalendarSuite) TearDownAllSuite() {
}

func (s *CalendarSuite) SetupTest() {
}

func (s *CalendarSuite) TestCreateEvent() {
	time.Sleep(time.Second * 10)
	client := http.DefaultClient

	createEditRequest := &internalhttp.CreateEditRequest{
		Title:       "event1",
		Time:        time.Now(),
		Duration:    10,
		Description: "event1 integrtest",
		UserID:      1,
		NotifyTime:  time.Now(),
	}

	jsonBody, _ := json.Marshal(createEditRequest)
	req, err := http.NewRequest(
		http.MethodPost, "http://app:8080/create",
		bytes.NewBuffer(jsonBody))

	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("error do request (%v)", err)
		return
	}

	timeEventRequest := &internalhttp.TimeEventRequest{
		Time: time.Now(),
	}

	jsonBody, _ = json.Marshal(timeEventRequest)
	req, err = http.NewRequest(
		http.MethodPost, "http://app:8080/day",
		bytes.NewBuffer(jsonBody))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	timeEventResponse := &internalhttp.TimeEventResponse{}

	_ = json.NewDecoder(resp.Body).Decode(timeEventResponse)

	if len(timeEventResponse.Events) != 1 {
		fmt.Printf("events count %v", len(timeEventResponse.Events))
		os.Exit(1)
	}

	os.Exit(0)
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}
