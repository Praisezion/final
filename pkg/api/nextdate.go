package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Vichislenie sled. dati
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("некорректный формат даты")
	}

	if repeat == "" {
		return "", nil
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("некорректный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 {
			return "", errors.New("некорректное количество дней")
		}
		if days > 400 {
			return "", errors.New("максимальное количество дней - 400")
		}

		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateFormat), nil
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "w":
		return handleWeekly(now, date, parts)

	case "m":
		return handleMonthly(now, date, parts)

	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

// Ezhenedelnoe povtorenie
func handleWeekly(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("некорректный формат правила w")
	}

	// Dni nedeli ***
	weekDays := make(map[time.Weekday]bool)
	daysStr := strings.Split(parts[1], ",")
	for _, d := range daysStr {
		day, err := strconv.Atoi(d)
		if err != nil || day < 1 || day > 7 {
			return "", errors.New("некорректный день недели")
		}
		var weekday time.Weekday
		switch day {
		case 1:
			weekday = time.Monday
		case 2:
			weekday = time.Tuesday
		case 3:
			weekday = time.Wednesday
		case 4:
			weekday = time.Thursday
		case 5:
			weekday = time.Friday
		case 6:
			weekday = time.Saturday
		case 7:
			weekday = time.Sunday
		}
		weekDays[weekday] = true
	}

	date = date.AddDate(0, 0, 1)

	for {

		if weekDays[date.Weekday()] && date.After(now) {
			break
		}
		date = date.AddDate(0, 0, 1)
	}
	return date.Format(DateFormat), nil
}

// Ezhemesyachnoe
func handleMonthly(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", errors.New("некорректный формат правила m")
	}

	days := make(map[int]bool)
	daysStr := strings.Split(parts[1], ",")
	for _, d := range daysStr {
		day, err := strconv.Atoi(d)
		if err != nil {
			return "", errors.New("некорректный день месяца")
		}
		if day < -2 || day > 31 || day == 0 {
			return "", errors.New("некорректный день месяца")
		}
		days[day] = true
	}

	months := make(map[time.Month]bool)
	if len(parts) == 3 {
		monthsStr := strings.Split(parts[2], ",")
		for _, m := range monthsStr {
			month, err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				return "", errors.New("некорректный месяц")
			}
			months[time.Month(month)] = true
		}
	} else {
		for m := 1; m <= 12; m++ {
			months[time.Month(m)] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)

		if !months[date.Month()] {
			continue
		}

		day := date.Day()
		lastDay := lastDayOfMonth(date)

		matched := false
		for d := range days {
			if d == day {
				matched = true
				break
			}
			if d == -1 && day == lastDay {
				matched = true
				break
			}
			if d == -2 && day == lastDay-1 {
				matched = true
				break
			}
		}

		if matched && date.After(now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
