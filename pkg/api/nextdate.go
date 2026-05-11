package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scheduler/pkg/db"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	repeatRule := strings.Fields(strings.TrimSpace(repeat))
	if len(repeatRule) == 0 {
		return "", fmt.Errorf("repeat rule is empty")
	}

	repeatType := repeatRule[0]

	date, err := time.Parse(db.DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format %q: %w", dstart, err)
	}

	var days int
	var months int
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

	case "w":
		if len(repeatRule) != 2 {
			return "", fmt.Errorf("repeat rule 'w' requires a list of weekdays")
		}

		parts := strings.Split(repeatRule[1], ",")
		weekdayNums := make(map[int]bool, len(parts))

		for _, p := range parts {
			n, convErr := strconv.Atoi(strings.TrimSpace(p))
			if convErr != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("repeat rule 'w' weekday must be a number between 1 and 7, got %q", p)
			}
			weekdayNums[n] = true
		}

		currentweekday := func(t time.Time) int {
			wd := int(t.Weekday())
			if wd == 0 {
				return 7
			}
			return wd
		}

		for {
			if date.After(now) && weekdayNums[currentweekday(date)] {
				return date.Format(db.DateFormat), nil
			}
			date = date.AddDate(0, 0, 1)
		}

	case "m":
		if len(repeatRule) < 2 || len(repeatRule) > 3 {
			return "", fmt.Errorf("repeat rule 'm' requires day list and optional month list (optional)")
		}

		dayParts := strings.Split(repeatRule[1], ",")
		dayNums := make([]int, 0, len(dayParts))
		for _, p := range dayParts {
			n, convErr := strconv.Atoi(strings.TrimSpace(p))
			if convErr != nil || (n < -2 || n == 0 || n > 31) {
				return "", fmt.Errorf("repeat rule 'm' day must be 1..31, -1 or -2, got %q", p)
			}
			dayNums = append(dayNums, n)
		}

		monthNums := make([]int, 0, 12)
		if len(repeatRule) == 3 {
			monthParts := strings.Split(repeatRule[2], ",")
			for _, p := range monthParts {
				n, convErr := strconv.Atoi(strings.TrimSpace(p))
				if convErr != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("repeat rule 'm' month must be 1..12, got %q", p)
				}
				monthNums = append(monthNums, n)
			}
		}

		monthAllowed := func(m time.Month) bool {
			if len(monthNums) == 0 {
				return true
			}
			for _, mn := range monthNums {
				if m == time.Month(mn) {
					return true
				}
			}
			return false
		}

		resolveDay := func(ruleDay, year int, month time.Month) (int, bool) {
			lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
			switch ruleDay {
			case -1:
				return lastDay, true
			case -2:
				if lastDay >= 2 {
					return lastDay - 1, true
				}
				return 0, false
			default:
				if ruleDay >= 1 && ruleDay <= lastDay {
					return ruleDay, true
				}
				return 0, false
			}
		}

		searchFrom := now
		if date.After(searchFrom) {
			searchFrom = date
		}

		currentMonth := time.Date(searchFrom.Year(), searchFrom.Month(), 1, 0, 0, 0, 0, time.UTC)
		for i := 0; i < 2400; i++ {
			year, month, _ := currentMonth.Date()
			if monthAllowed(month) {
				var best *time.Time
				for _, ruleDay := range dayNums {
					day, ok := resolveDay(ruleDay, year, month)
					if !ok {
						continue
					}
					candidate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
					if !candidate.After(searchFrom) {
						continue
					}
					if best == nil || candidate.Before(*best) {
						c := candidate
						best = &c
					}
				}
				if best != nil {
					return best.Format(db.DateFormat), nil
				}
			}
			currentMonth = currentMonth.AddDate(0, 1, 0)
		}
		return "", fmt.Errorf("repeat rule 'm': no matching date found")

	default:
		return "", fmt.Errorf("repeat rule has unsupported type: %q", repeatType)
	}

	for {
		date = date.AddDate(years, months, days)
		if date.After(now) {
			break
		}
	}

	return date.Format(db.DateFormat), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {

		now, err = time.Parse(db.DateFormat, nowStr)
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
