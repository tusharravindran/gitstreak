package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Config struct {
	Username     string `json:"username"`
	ReminderHour int    `json:"reminder_hour"`
	ReminderMin  int    `json:"reminder_minute"`
	SkipDays     []int  `json:"skip_days"` // 0=Sunday … 6=Saturday

	// LastAuditedFiles maps date (YYYY-MM-DD) to the lone file touched that day, when
	// that day's contribution was a single-file commit. Used to detect the same file
	// being farmed for streak days in a row. Pruned to the last 7 entries on save.
	LastAuditedFiles map[string]string `json:"last_audited_files,omitempty"`
}

func Default() Config {
	return Config{
		ReminderHour: 21,
		ReminderMin:  0,
		SkipDays:     []int{},
	}
}

func path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gitstreak", "config.json"), nil
}

func Load() (Config, error) {
	cfg := Default()
	p, err := path()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func Save(cfg Config) error {
	p, err := path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

func (c Config) ReminderLabel() string {
	ampm := "AM"
	h := c.ReminderHour
	if h >= 12 {
		ampm = "PM"
		if h > 12 {
			h -= 12
		}
	}
	if h == 0 {
		h = 12
	}
	return fmt.Sprintf("%d:%02d %s", h, c.ReminderMin, ampm)
}

// IsTodaySkipped returns true if today is in the skip list
func (c Config) IsTodaySkipped() bool {
	today := int(time.Now().Weekday())
	for _, d := range c.SkipDays {
		if d == today {
			return true
		}
	}
	return false
}

// ConsecutiveSkipCount returns how many consecutive days starting from Sunday are skipped
// (used to detect >3 consecutive skip days)
func ConsecutiveSkipCount(days []int) int {
	if len(days) == 0 {
		return 0
	}
	set := map[int]bool{}
	for _, d := range days {
		set[d] = true
	}
	max := 0
	run := 0
	// check 0–6 twice to handle wrap-around (Sun..Sat..Sun)
	for i := 0; i < 14; i++ {
		if set[i%7] {
			run++
			if run > max {
				max = run
			}
		} else {
			run = 0
		}
	}
	return max
}

// ParseDays parses comma-separated day names into weekday ints (0=Sun…6=Sat)
func ParseDays(input string) ([]int, error) {
	names := map[string]int{
		"sun": 0, "sunday": 0,
		"mon": 1, "monday": 1,
		"tue": 2, "tuesday": 2,
		"wed": 3, "wednesday": 3,
		"thu": 4, "thursday": 4,
		"fri": 5, "friday": 5,
		"sat": 6, "saturday": 6,
	}
	var result []int
	seen := map[int]bool{}
	for _, part := range strings.Split(input, ",") {
		key := strings.ToLower(strings.TrimSpace(part))
		if key == "" {
			continue
		}
		v, ok := names[key]
		if !ok {
			return nil, fmt.Errorf("unknown day %q — use Mon, Tue, Wed, Thu, Fri, Sat, Sun", part)
		}
		if !seen[v] {
			result = append(result, v)
			seen[v] = true
		}
	}
	return result, nil
}

// DayName returns the weekday name for a given int
func DayName(d int) string {
	return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}[d%7]
}

// RecordAuditedFile stores the lone file touched on the given date, pruning history
// to the last 7 entries. Pass "" to clear the entry for a date (e.g. multi-file commit).
func (c *Config) RecordAuditedFile(date, file string) {
	if c.LastAuditedFiles == nil {
		c.LastAuditedFiles = map[string]string{}
	}
	if file == "" {
		delete(c.LastAuditedFiles, date)
	} else {
		c.LastAuditedFiles[date] = file
	}

	if len(c.LastAuditedFiles) <= 7 {
		return
	}
	dates := make([]string, 0, len(c.LastAuditedFiles))
	for d := range c.LastAuditedFiles {
		dates = append(dates, d)
	}
	sort.Strings(dates)
	for _, d := range dates[:len(dates)-7] {
		delete(c.LastAuditedFiles, d)
	}
}

// RecentSingleFiles returns the lone-file history for the n days immediately before
// (not including) the given date, oldest first — used to detect "same file" farming.
func (c Config) RecentSingleFiles(beforeDate string, n int) []string {
	before, err := time.Parse("2006-01-02", beforeDate)
	if err != nil {
		return nil
	}
	var out []string
	for i := n; i >= 1; i-- {
		d := before.AddDate(0, 0, -i).Format("2006-01-02")
		if f, ok := c.LastAuditedFiles[d]; ok {
			out = append(out, f)
		} else {
			out = append(out, "")
		}
	}
	return out
}

// SkipDayNames returns human-readable skip day names
func (c Config) SkipDayNames() []string {
	var names []string
	for _, d := range c.SkipDays {
		names = append(names, DayName(d))
	}
	return names
}
