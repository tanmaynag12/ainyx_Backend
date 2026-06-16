package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "birthday already passed this year",
			dob:      time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		{
			name:     "birthday not yet this year",
			dob:      time.Date(1995, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		{
			name:     "born today same month same day",
			dob:      time.Date(2003, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
			expected: 23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateAge(tt.dob)
			if got != tt.expected {
				t.Errorf("calculateAge(%v) = %d, want %d", tt.dob, got, tt.expected)
			}
		})
	}
}

func TestDaysUntilNextBirthday(t *testing.T) {
	today := time.Now()

	nextYear := time.Date(today.Year()+1, today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	expectedDays := int(nextYear.Sub(today).Hours() / 24)

	dob := time.Date(1995, today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	got := daysUntilNextBirthday(dob)

	if got != expectedDays {
		t.Errorf("daysUntilNextBirthday() = %d, want %d", got, expectedDays)
	}
}
