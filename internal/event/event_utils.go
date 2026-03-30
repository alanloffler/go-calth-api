package event

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const timeZoneOffset = -3

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

func localLoc() *time.Location {
	return time.FixedZone("ART", timeZoneOffset*3600)
}

func isWithinSchedule(candidate time.Time, profile sqlc.ProfessionalProfile) bool {
	loc := localLoc()
	localTime := candidate.In(loc)

	dayOfWeek := strconv.Itoa(int(localTime.Weekday()))
	workingDays := strings.Split(profile.WorkingDays, ",")
	dayOk := false

	for _, d := range workingDays {
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
	totalMinutes := localTime.Hour()*60 + localTime.Minute()

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

func areAllSlotsFree(
	ctx context.Context,
	repo *EventRepository,
	candidateStart time.Time,
	days int,
	slotDuration int,
	profile sqlc.ProfessionalProfile,
	businessID, professionalID pgtype.UUID) (bool, error) {
	if !isWithinSchedule(candidateStart, profile) {
		return false, nil
	}

	slotDur := time.Duration(slotDuration) * time.Minute
	candidateDates := generateRecurringDates(candidateStart, int32(days))

	for _, cd := range candidateDates {
		if !isWithinSchedule(cd, profile) {
			return false, nil
		}

		slotEnd := cd.Add(slotDur)
		_, err := repo.ChechSlotConflict(ctx, sqlc.CheckSlotConflictParams{
			BusinessID:     businessID,
			ProfessionalID: professionalID,
			EndDate:        pgtype.Timestamptz{Time: cd, Valid: true},
			StartDate:      pgtype.Timestamptz{Time: slotEnd, Valid: true},
		})
		if err == nil {
			return false, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return false, err
		}
	}

	return true, nil
}

func findSuggestion(
	ctx context.Context,
	repo *EventRepository,
	startDate time.Time,
	days int,
	profile sqlc.ProfessionalProfile,
	businessID, professionalID pgtype.UUID) (*time.Time, error) {
	slotDuration, err := strconv.Atoi(profile.SlotDuration)
	if err != nil {
		return nil, err
	}

	startH, startM := parseHourMinute(profile.StartHour)
	endH, endM := parseHourMinute(profile.EndHour)
	scheduleStart := startH*60 + startM
	scheduleEnd := endH*60 + endM

	// Same day, other slots
	loc := localLoc()
	localStart := startDate.In(loc)

	for min := scheduleStart; min+slotDuration <= scheduleEnd; min += slotDuration {
		localCandidate := time.Date(
			localStart.Year(), localStart.Month(), localStart.Day(),
			min/60, min%60, 0, 0, loc,
		)
		candidate := localCandidate.UTC()
		if candidate.Equal(startDate) {
			continue
		}
		free, err := areAllSlotsFree(ctx, repo, candidate, days, slotDuration, profile, businessID, professionalID)
		if err != nil {
			return nil, err
		}
		if free {
			return &candidate, nil
		}
	}

	// Next 7 days, same hour
	for day := 1; day <= 7; day++ {
		candidate := startDate.AddDate(0, 0, day)
		free, err := areAllSlotsFree(ctx, repo, candidate, days, slotDuration, profile, businessID, professionalID)
		if err != nil {
			return nil, err
		}
		if free {
			return &candidate, nil
		}
	}

	return nil, nil
}
