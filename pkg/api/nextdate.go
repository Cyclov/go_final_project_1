package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	repeatRule := strings.Split(strings.TrimSpace(repeat), " ") // trim spaces and split by space

	repeatType := repeatRule[0]

	// Fix: time.Parse returns (time.Time, error) — both values must be captured
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format %q: %w", dstart, err)
	}

	var days int
	var years int

	switch repeatType {
	case "d":
		if len(repeatRule) != 2 {
			return "", fmt.Errorf("repeat rule 'd' requires one interval parameter")
		}
		days, err = strconv.Atoi(repeatRule[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("repeat rule 'd' interval must be a number between 1 and 400")
		}

	case "y":
		years = 1

	//case "w":
	//case "m":
	default:
		return "", fmt.Errorf("repeat rule has unsupported type: %q", repeatType)
	}

	for {
		date = date.AddDate(years, 0, days)
		if date.After(now) {
			break
		}
	}

	return date.Format("20060102"), nil
}

func InitNextDayHandler() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {

		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid 'now' date: %v", err), http.StatusBadRequest)
			return
		}
	}

	dateStart := r.FormValue("date")
	if dateStart == "" {
		http.Error(w, "missing 'date start' parameter", http.StatusBadRequest)
		return
	}

	repeat := r.FormValue("repeat")
	if repeat == "" {
		http.Error(w, "missing 'repeat' parameter", http.StatusBadRequest)
		return
	}

	result, err := NextDate(now, dateStart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, result)
}
