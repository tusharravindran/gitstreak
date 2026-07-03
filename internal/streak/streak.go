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

func Calculate(days []gh.ContributionDay) Result {
	if len(days) == 0 {
		return Result{}
	}

	today := time.Now().Format("2006-01-02")

	var totalThisYear int
	for _, d := range days {
		totalThisYear += d.ContributionCount
	}

	// today's count
	var todayCount int
	for _, d := range days {
		if d.Date == today {
			todayCount = d.ContributionCount
			break
		}
	}

	// current streak — walk backwards from today
	current := 0
	for i := len(days) - 1; i >= 0; i-- {
		d := days[i]
		if d.Date > today {
			continue
		}
		// allow today even with 0 (streak is still alive)
		if d.Date == today {
			if todayCount > 0 {
				current++
			}
			continue
		}
		if d.ContributionCount == 0 {
			break
		}
		current++
	}

	// longest streak
	longest := 0
	run := 0
	for _, d := range days {
		if d.ContributionCount > 0 {
			run++
			if run > longest {
				longest = run
			}
		} else {
			run = 0
		}
	}

	// last active date
	lastActive := ""
	for i := len(days) - 1; i >= 0; i-- {
		if days[i].ContributionCount > 0 {
			lastActive = days[i].Date
			break
		}
	}

	return Result{
		CurrentStreak:  current,
		LongestStreak:  longest,
		TodayCount:     todayCount,
		TotalThisYear:  totalThisYear,
		CommittedToday: todayCount > 0,
		LastActiveDate: lastActive,
		HeatMap:        days,
	}
}
