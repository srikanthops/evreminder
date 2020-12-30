package main

import (
	"fmt"
	"net/http"

	ev "github.com/srikanthops/evreminder/event"
)

func todayEvents(w http.ResponseWriter, r *http.Request) {
	for _, e := range ev.TodayEvents() {
		fmt.Fprintf(w, ev.EvStr(e))
	}
}

func tomEvents(w http.ResponseWriter, r *http.Request) {
	for _, e := range ev.TomEvents() {
		fmt.Fprintf(w, ev.EvStr(e))
	}
}

func yEvents(w http.ResponseWriter, r *http.Request) {
	for _, e := range ev.YEvents() {
		fmt.Fprintf(w, ev.EvStr(e))
	}
}

func mEvents(w http.ResponseWriter, r *http.Request) {
	for _, e := range ev.MEvents() {
		fmt.Fprintf(w, ev.EvStr(e))
	}
}
