package domain

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func IsValidTitle(t string) bool {
	l := utf8.RuneCountInString(strings.TrimSpace(t))
	if l < 3 || l > 120 {
		return false
	}

	return true
}

func IsValidDesc(d string) bool {
	return utf8.RuneCountInString(d) <= 200
}

func IsValidYear(y string) bool {
	year, err := strconv.Atoi(y)
	if err != nil {
		return false
	}

	if year < 2026 || year > time.Now().Year() {
		return false
	}

	return true
}

func IsValidDate(d string) bool {
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		return false
	}

	if date.Year() < 2026 {
		return false
	}

	now := time.Now()

	limit := now.Add(24 * time.Hour)

	if date.After(limit) {
		return false
	}

	return true
}
