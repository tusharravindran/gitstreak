package streak

import (
	"time"

	gh "github.com/tusharravindran/gitstreak/internal/github"
)

type Result struct {
	CurrentStreak  int
	LongestStreak  int
	TodayCount     int
	TotalThisYear  int
	CommittedToday bool
	LastActiveDate string
	HeatMap        []gh.ContributionDay
}

// isSkipDay returns true if the given date falls on a skip weekday
func isSkipDay(date string, skipDays []int) bool {
	if len(skipDays) == 0 {
		return false
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	wd := int(t.Weekday())
	for _, d := range skipDays {
		if d == wd {
			return true
		}
	}
	return false
}

func Calculate(days []gh.ContributionDay, skipDays []int) Result {
	if len(days) == 0 {
		return Result{}
	}

	today := time.Now().Format("2006-01-02")

	var totalThisYear int
	for _, d := range days {
		totalThisYear += d.ContributionCount
	}

	var todayCount int
	for _, d := range days {
		if d.Date == today {
			todayCount = d.ContributionCount
			break
		}
	}

	// current streak — walk backwards from today
	// skip days with 0 contributions are treated as non-breaking if they're skip days
	current := 0
	for i := len(days) - 1; i >= 0; i-- {
		d := days[i]
		if d.Date > today {
			continue
		}
		if d.Date == today {
			// today is a skip day with no commit — streak still alive, just don't count it
			if isSkipDay(d.Date, skipDays) && todayCount == 0 {
				continue
			}
			if todayCount > 0 {
				current++
			}
			continue
		}
		// past days: skip days with 0 contributions don't break the streak
		if d.ContributionCount == 0 {
			if isSkipDay(d.Date, skipDays) {
				continue // transparent — doesn't break or add to streak
			}
			break
		}
		current++
	}

	// longest streak (also respecting skip days)
	longest := 0
	run := 0
	for _, d := range days {
		if d.ContributionCount > 0 {
			run++
			if run > longest {
				longest = run
			}
		} else if isSkipDay(d.Date, skipDays) {
			// skip day with no commit — don't break the run, don't add to it
			continue
		} else {
			run = 0
		}
	}

	lastActive := ""
	for i := len(days) - 1; i >= 0; i-- {
		if days[i].ContributionCount > 0 {
			lastActive = days[i].Date
			break
		}
	}

	// CommittedToday is true if today is a skip day (nothing required) or has commits
	todayIsSkip := isSkipDay(today, skipDays)
	committedToday := todayCount > 0 || todayIsSkip

	return Result{
		CurrentStreak:  current,
		LongestStreak:  longest,
		TodayCount:     todayCount,
		TotalThisYear:  totalThisYear,
		CommittedToday: committedToday,
		LastActiveDate: lastActive,
		HeatMap:        days,
	}
}
