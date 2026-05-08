package event

import (
	"strconv"
	"strings"
	"time"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

const timeZoneOffset = -3

type interval struct {
	start int64
	end   int64
}

type RecurringSuggestions struct {
	SameDay           []time.Time `json:"sameDay"`
	OtherDaysSameHour []time.Time `json:"otherDaysSameHour"`
	OtherDaysAnyHour  []time.Time `json:"otherDaysAnyHour"`
}

func localLoc() *time.Location {
	return time.FixedZone("ART", timeZoneOffset*3600)
}

func generateRecurringDates(start time.Time, occurrences int32) []time.Time {
	dates := make([]time.Time, 0, occurrences)
	for i := int32(0); i < occurrences; i++ {
		dates = append(dates, start.AddDate(0, 0, int(i)*7))
	}
	return dates
}

func parseHourMinute(s string) (int, int) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, 0
	}
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return h, m
}

func isWithinSchedule(candidate time.Time, profile sqlc.ProfessionalProfile) bool {
	local := candidate.In(localLoc())

	dayOfWeek := strconv.Itoa(int(local.Weekday()))
	dayOk := false
	for _, d := range strings.Split(profile.WorkingDays, ",") {
		if strings.TrimSpace(d) == dayOfWeek {
			dayOk = true
			break
		}
	}
	if !dayOk {
		return false
	}

	slotDuration, err := strconv.Atoi(profile.SlotDuration)
	if err != nil {
		return false
	}

	startH, startM := parseHourMinute(profile.StartHour)
	endH, endM := parseHourMinute(profile.EndHour)
	scheduleStart := startH*60 + startM
	scheduleEnd := endH*60 + endM
	totalMinutes := local.Hour()*60 + local.Minute()

	if totalMinutes < scheduleStart || totalMinutes+slotDuration > scheduleEnd {
		return false
	}

	if profile.DailyExceptionStart.Valid && profile.DailyExceptionEnd.Valid {
		exH, exM := parseHourMinute(profile.DailyExceptionStart.String)
		exEndH, exEndM := parseHourMinute(profile.DailyExceptionEnd.String)
		if totalMinutes >= exH*60+exM && totalMinutes < exEndH*60+exEndM {
			return false
		}
	}

	return true
}

func conflictsExisting(intervals []interval, start, end int64) bool {
	for _, iv := range intervals {
		if iv.start >= end {
			break
		}
		if iv.end > start {
			return true
		}
	}
	return false
}

func buildRecurringSuggestions(
	originalStart time.Time,
	days int,
	profile sqlc.ProfessionalProfile,
	intervals []interval,
	slotMs int64,
) RecurringSuggestions {
	const maxPerTier = 3

	slotDurationMin, _ := strconv.Atoi(profile.SlotDuration)
	startH, startM := parseHourMinute(profile.StartHour)
	endH, endM := parseHourMinute(profile.EndHour)
	scheduleStart := startH*60 + startM
	scheduleEnd := endH*60 + endM

	isSeriesFree := func(candidate time.Time) bool {
		if !isWithinSchedule(candidate, profile) {
			return false
		}
		for _, c := range generateRecurringDates(candidate, int32(days)) {
			if !isWithinSchedule(c, profile) {
				return false
			}
			startMs := c.UnixMilli()
			if conflictsExisting(intervals, startMs, startMs+slotMs) {
				return false
			}
		}
		return true
	}

	loc := localLoc()
	localOrig := originalStart.In(loc)
	origY, origMo, origD := localOrig.Date()

	// Tier 1: same local day, different hours.
	sameDay := make([]time.Time, 0, maxPerTier)
	for min := scheduleStart; min+slotDurationMin <= scheduleEnd && len(sameDay) < maxPerTier; min += slotDurationMin {
		candidate := time.Date(origY, origMo, origD, min/60, min%60, 0, 0, loc).UTC()
		if candidate.Equal(originalStart) {
			continue
		}
		if isSeriesFree(candidate) {
			sameDay = append(sameDay, candidate)
		}
	}

	// Tier 2: other days, same hour.
	otherDaysSameHour := make([]time.Time, 0, maxPerTier)
	for day := 1; day <= 7 && len(otherDaysSameHour) < maxPerTier; day++ {
		candidate := localOrig.AddDate(0, 0, day).UTC()
		if isSeriesFree(candidate) {
			otherDaysSameHour = append(otherDaysSameHour, candidate)
		}
	}

	// Tier 3: other days, any hour, deduped against Tier 2.
	tier2 := make(map[int64]struct{}, len(otherDaysSameHour))
	for _, t := range otherDaysSameHour {
		tier2[t.UnixMilli()] = struct{}{}
	}
	otherDaysAnyHour := make([]time.Time, 0, maxPerTier)
outer:
	for day := 1; day <= 7; day++ {
		shifted := localOrig.AddDate(0, 0, day)
		y, mo, d := shifted.Date()
		for min := scheduleStart; min+slotDurationMin <= scheduleEnd; min += slotDurationMin {
			if len(otherDaysAnyHour) >= maxPerTier {
				break outer
			}
			candidate := time.Date(y, mo, d, min/60, min%60, 0, 0, loc).UTC()
			if _, dup := tier2[candidate.UnixMilli()]; dup {
				continue
			}
			if isSeriesFree(candidate) {
				otherDaysAnyHour = append(otherDaysAnyHour, candidate)
			}
		}
	}

	return RecurringSuggestions{
		SameDay:           sameDay,
		OtherDaysSameHour: otherDaysSameHour,
		OtherDaysAnyHour:  otherDaysAnyHour,
	}
}
