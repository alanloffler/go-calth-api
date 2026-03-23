package event

import "time"

func generateRecurringDates(start time.Time, occurrences int32) []time.Time {
	dates := make([]time.Time, 0, occurrences)

	for i := int32(0); i < occurrences; i++ {
		dates = append(dates, start.AddDate(0, 0, int(i)*7))
	}

	return dates
}
