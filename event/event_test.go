package event

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var input = `
[
{
    "name": "JanBirthday", 
    "date": "2018-01-12 04:35",
    "evType": "Birthday"
},
{
    "name": "FebBirthday", 
    "date": "2018-02-11 04:35",
    "evType": "Birthday"
},
{
    "name": "MarchBirthday", 
    "date": "2018-03-10 08:35",
    "evType": "Birthday"
},
{
    "name": "AprilBirthday", 
    "date": "2019-04-06 09:35",
    "evType": "Birthday"
},
{
    "name": "christmas", 
    "date": "2011-12-15 09:35",
    "evType": "Birthday"
},
{
    "name": "GBirthday", 
    "date": "1972-12-06 09:35",
    "evType": "Birthday"
},
{
    "name": "lDeathDay", 
    "date": "2019-12-10 09:35",
    "evType": "Death"
}
]
`

const (
	oneDay = time.Duration(time.Hour * 24)
)

var testEvents Events
var buffEvents *Events

func addTestEvent(evName string, evTime time.Time, evType EvType) {
	ev := AnnualEvent{evName,
		evTime,
		evType,
	}
	testEvents.evCount++
	testEvents.events = append(testEvents.events, ev)

}
func init() {
	now := time.Now()

	after := now.AddDate(1, 0, 0)
	fmt.Println("\nAdd 1 Year:", after)

	after = now.AddDate(0, 1, 0)
	fmt.Println("\nAdd 1 Month:", after)

	after = now.AddDate(0, 0, -1)
	fmt.Println("\nAdd 1 Day:", after)

	after = now.AddDate(2, 2, 5)
	fmt.Println("\nAdd multiple values:", after)

	addTestEvent("todayBirthday", time.Now(), EventBirthday)
	addTestEvent("todayAnniversary", time.Now(), EventAnniversary)
	addTestEvent("yAnniversary", time.Now().AddDate(-10, 0, -1), EventAnniversary)
	addTestEvent("yBirthday", time.Now().AddDate(-2, 0, -1), EventBirthday)
	addTestEvent("tomBirthday", time.Now().AddDate(-2, 0, 1), EventBirthday)
	addTestEvent("tomAnniversary", time.Now().AddDate(-3, 0, 1), EventAnniversary)
	addTestEvent("tomDeath", time.Now().Add(oneDay), EventDeath)

	fmt.Fprintf(os.Stdout, "evCount=%d\n", testEvents.evCount)
	for _, ev := range testEvents.events {
		fmt.Fprintf(os.Stdout, "ev:%s %v\n", ev.name, ev.date)
	}
}

func TestTodayEvents(t *testing.T) {
	count := 0
	expCount := 2
	for _, ev := range testEvents.events {
		if isEventToday(ev) {
			count++
		}
	}
	if count != expCount {
		t.Errorf("Fail Today expected count :%d got :%d", expCount, count)
	}
}

func TestYesterdayEvents(t *testing.T) {
	count := 0
	expCount := 2
	for _, ev := range testEvents.events {
		if isEventYesterday(ev) {
			count++
		}
	}
	if count != expCount {
		t.Errorf("Fail Yesterday expected count :%d got :%d", expCount, count)
	}

}

func TestTomorowEvents(t *testing.T) {
	count := 0
	expCount := 3
	for _, ev := range testEvents.events {
		if isEventTomorrow(ev) {
			count++
		}
	}
	if count != expCount {
		t.Errorf("Fail Tomorrow expected count :%d got :%d", expCount, count)
	}
}

func TestThisMonth(t *testing.T) {
	count := 0
	expCount := 7
	for _, ev := range testEvents.events {
		if isEventThisMonth(ev) {
			count++
		}
	}
	if count != expCount {
		t.Errorf("Fail This Month expected count :%d got :%d", expCount, count)
	}
}

func TestBuffParse(t *testing.T) {
	evCExpected := 7
	evs, err := populateEventsFromBuf([]byte(input))

	if err != nil {
		t.Errorf("Error parsing input buffer:%v\n", err)
	}

	if evs != nil && evs.evCount != evCExpected {
		t.Errorf("Event count expected %d got:%d\n", evCExpected, evs.evCount)
	}

	for _, e := range evs.events {
		fmt.Fprintf(os.Stdout, " ev:%s time:%s\n", e.name, e.date)
	}
	buffEvents = evs
}

// TODO This test will only work in December
func TestBufferDec(t *testing.T) {
	count := 0
	expCount := 3
	for _, ev := range buffEvents.events {
		if isEventThisMonth(ev) {
			count++
		}
	}
	if count != expCount {
		t.Errorf("Fail This Month expected count :%d got :%d", expCount, count)
	}
}
