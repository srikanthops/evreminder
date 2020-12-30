package event

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// EvType Event Type
type EvType string

const (
	// EventBirthday Birthday Event
	EventBirthday EvType = "Birthday"
	// EventAnniversary Event
	EventAnniversary EvType = "Anniversary"
	// EventDeath Event
	EventDeath EvType = "Death"
	// EventOther  Event
	EventOther EvType = "Other"
)

// AnnualEvent Annual Events Struct
type AnnualEvent struct {
	name   string
	date   time.Time
	evType EvType
}

// EvStr formatted event
func EvStr(ev AnnualEvent) string {
	return fmt.Sprintf("%s %s %s\n", ev.name, ev.date.Format("2006-01-02"), ev.evType)
}

// Events data struct
type Events struct {
	events  []AnnualEvent
	evCount int
	source  string
}

// EvImpl Events Implementation
type EvImpl struct{}

//EVHandler global handler
var EVHandler EvImpl

var evc *Events

// How do we populate events from google bucket or static file

// PopulateEvents populat events from a file
func PopulateEvents(f string) error {
	// Open our jsonFile
	jsonFile, err := os.Open(f)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	evc, err = populateEventsFromBuf(byteValue)
	return err
}

func populateEventsFromBuf(buf []byte) (*Events, error) {

	type jsEvent struct {
		EVname string `json:"name"`
		EVDate string `json:"date"`
		EVType string `json:"evType"`
	}

	var ev = Events{}
	evc = &ev
	//var evs  []AnnualEvent

	var evs []jsEvent

	if err := json.Unmarshal([]byte(buf), &evs); err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stdout, ">> evLen=%d\n", len(evs))

	for _, jev := range evs {
		fmt.Fprintf(os.Stdout, "ev:%s time:%s\n", jev.EVname, jev.EVDate)
		d, err := time.Parse("2006-01-02 15:04", jev.EVDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erroring pasring date %s", jev.EVDate)
			continue
		}

		ev.events = append(ev.events, AnnualEvent{jev.EVname, d, EvType(jev.EVType)})
		ev.evCount++
	}
	ev.source = "buffer"
	return &ev, nil
}

func (ev *AnnualEvent) day() int {
	_, _, d := ev.date.Date()
	return d
}

func (ev *AnnualEvent) month() time.Month {
	_, m, _ := ev.date.Date()
	return m
}

func (ev *AnnualEvent) year() int {
	y, _, _ := ev.date.Date()
	return y
}

// Return True if today's Date and Month Match
func isEventToday(ev AnnualEvent) bool {
	_, m, d := time.Now().Date()
	if ev.month() == m && ev.day() == d {
		return true
	}
	return false
}

// Return True if AnnualEvent is yesterday
func isEventYesterday(ev AnnualEvent) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	_, m, d := yesterday.Date()
	if ev.month() == m && ev.day() == d {
		return true
	}
	return false
}

// Return True if AnnualEvent is tomorrow
func isEventTomorrow(ev AnnualEvent) bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	_, m, d := tomorrow.Date()
	if ev.month() == m && ev.day() == d {
		return true
	}
	return false
}

// Return True if AnnualEvent is this month
func isEventThisMonth(ev AnnualEvent) bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	_, m, _ := tomorrow.Date()
	if ev.month() == m {
		return true
	}
	return false
}

// Return True if the AnnualEvent belongs to this week.
// Week begins with Monday ->  and ends with Sunday
func isEventThisWeek(ev AnnualEvent) bool {
	// TBD
	return false
}

type filter func(AnnualEvent) bool

func eventsFilter(f filter) []AnnualEvent {
	evs := []AnnualEvent{}
	for _, ev := range evc.events {
		if f(ev) {
			evs = append(evs, ev)
		}
	}
	return evs
}

// TodayEvents today events
func TodayEvents() []AnnualEvent {
	return eventsFilter(isEventToday)
}

// TomEvents tomorrow events
func TomEvents() []AnnualEvent {
	return eventsFilter(isEventTomorrow)
}

// YEvents yesterday events
func YEvents() []AnnualEvent {
	return eventsFilter(isEventYesterday)
}

// MEvents Month events
func MEvents() []AnnualEvent {
	return eventsFilter(isEventThisMonth)
}
