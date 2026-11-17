package domain_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/N1ktarchik/habit-graph/internal/core/domain"
)

func TestIsValidTitle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "valid standard title", input: "Read Go books", expected: true},
		{name: "valid min length (3 runes)", input: "Run", expected: true},
		{name: "valid min length cyrillic (3 runes)", input: "Бег", expected: true},
		{name: "valid max length (120 runes)", input: strings.Repeat("а", 120), expected: true},
		{name: "valid with leading and trailing spaces", input: "   Gym   ", expected: true},
		{name: "invalid empty string", input: "", expected: false},
		{name: "invalid only spaces", input: "     ", expected: false},
		{name: "invalid too short (2 runes)", input: "Go", expected: false},
		{name: "invalid too short after trim", input: "  ab  ", expected: false},
		{name: "invalid too long (121 runes)", input: strings.Repeat("a", 121), expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, domain.IsValidTitle(tc.input))
		})
	}
}

func TestIsValidDesc(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "valid empty description", input: "", expected: true},
		{name: "valid normal description", input: "30 minutes every morning", expected: true},
		{name: "valid max length (200 runes cyrillic)", input: strings.Repeat("я", 200), expected: true},
		{name: "invalid too long (201 runes)", input: strings.Repeat("a", 201), expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, domain.IsValidDesc(tc.input))
		})
	}
}

func TestIsValidYear(t *testing.T) {
	t.Parallel()

	currentYear := time.Now().Year()

	testCases := []struct {
		name        string
		input       string
		expectNil   bool
		expectedVal int
	}{
		{name: "valid minimum year 2026", input: "2026", expectNil: false, expectedVal: 2026},
		{name: "valid current year", input: strconv.Itoa(currentYear), expectNil: false, expectedVal: currentYear},
		{name: "invalid year before 2026", input: "2025", expectNil: true},
		{name: "invalid future year", input: strconv.Itoa(currentYear + 1), expectNil: true},
		{name: "invalid non-numeric string", input: "abcd", expectNil: true},
		{name: "invalid empty string", input: "", expectNil: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := domain.IsValidYear(tc.input)
			if tc.expectNil {
				assert.Nil(t, res)
			} else {
				require.NotNil(t, res)
				assert.Equal(t, tc.expectedVal, *res)
			}
		})
	}
}

func TestIsValidDate(t *testing.T) {
	t.Parallel()

	today := time.Now().Format(time.DateOnly)
	tooFarFuture := time.Now().Add(48 * time.Hour).Format(time.DateOnly)

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "valid today date", input: today, expected: true},
		{name: "valid start of 2026", input: "2026-01-01", expected: true},
		{name: "invalid year 2025", input: "2025-12-31", expected: false},
		{name: "invalid future date beyond 24h limit", input: tooFarFuture, expected: false},
		{name: "invalid date format (DD-MM-YYYY)", input: "16-11-2026", expected: false},
		{name: "invalid non-existent calendar date", input: "2026-02-30", expected: false},
		{name: "invalid empty string", input: "", expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, domain.IsValidDate(tc.input))
		})
	}
}
