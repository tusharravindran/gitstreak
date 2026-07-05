package roast

import (
	"math/rand"
	"time"
)

var roasts = []string{
	"💀 Bro really said 'I'll commit tomorrow' for %d days straight.",
	"🪦 %d days offline. Your GitHub is starting to look like abandoned software.",
	"🌵 The tumbleweeds on your profile are getting lonely. %d days, really?",
	"😬 %d days without a commit. Your future self is filing a complaint.",
	"🧪 Scientists confirmed: %d days of no commits causes irreversible cringe.",
	"📞 GitHub support called. They thought your account was hacked — turns out you just ghosted for %d days.",
	"😶 %d days clean. But not the good kind of clean.",
	"👻 Your streak died %d days ago. Moment of silence.",
	"😳 Even your README is embarrassed. %d days and counting.",
	"🤡 %d days without a commit — bold strategy. How's that working out?",
	"📉 %d days of silence. Your contribution graph looks like a flatline.",
	"🫠 %d days. The compiler misses you. Actually, everything misses you.",
}

var skipRoasts = []string{
	"😬 Taking %d consecutive days off? Bold of you to assume your motivation survives that.",
	"🛋️ %d skip days in a row? That's not a schedule, that's a retirement plan.",
	"🏜️ Skipping %d days straight. Your commit graph is about to look like a desert.",
	"😅 %d consecutive off-days? The GitHub graph called — it barely recognizes you anymore.",
	"🥲 Careful — %d days away and you'll need to re-learn how to open your editor.",
	"😤 Setting %d skip days. Respect the honesty. Fear the consequences.",
}

var milestones = map[int]string{
	4:   "🔥 4 days straight. You're actually doing it. Keep going.",
	7:   "🏅 One full week! Consistency is a skill — you're building it.",
	14:  "⚡ Two weeks in. Habits are forming. This is where it gets real.",
	30:  "🏆 30 DAYS! A whole month of showing up. That's genuinely rare.",
	50:  "🔥🔥🔥 50 DAYS! Half a century of commits. You're built different.",
	75:  "💎 75 days. Most people quit at 30. You didn't. That says everything.",
	100: "🚀🚀 100 DAYS! Triple digits. Legendary. Absolutely legendary.",
	150: "🧠 150 days. At this point your keyboard is part of your identity.",
	200: "🏅✨ 200 days! GitHub should send you a trophy. This is elite-level.",
	365: "🌟🎉 A FULL YEAR! 365 days of showing up. You absolute unit. Incredible.",
}

var earlyPraise = []string{
	"✅ Committed today. That's how streaks start.",
	"🌱 Nice — another day, another commit.",
	"👏 Committed. Keep the chain going.",
}

var appreciations = []string{
	"🔥 %d days and still going. The momentum is real.",
	"⚡ %d-day streak. You're making it look easy.",
	"💪 %d days straight. Discipline is quietly compounding.",
	"🚀 %d days in. Stay consistent — compound interest kicks in soon.",
	"✨ %d days. Small daily actions, massive long-term results.",
	"😤 %d days. You're not streaking — you're building a habit.",
	"🎯 %d days. Most devs talk about side projects. You're shipping them.",
}

func rng() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func pick(list []string) string {
	return list[rng().Intn(len(list))]
}

// ForBrokenStreak returns a roast for n days without a commit
func ForBrokenStreak(daysMissed int) string {
	if daysMissed < 1 {
		return ""
	}
	return pick(roasts)
}

// ForSkipDays returns a roast for setting too many consecutive skip days
func ForSkipDays(consecutive int) string {
	return pick(skipRoasts)
}

// ForStreak returns appreciation or milestone message for an active streak
func ForStreak(streak int) string {
	// Check exact milestones first
	if msg, ok := milestones[streak]; ok {
		return msg
	}
	// Generic appreciation for 4+
	if streak >= 4 {
		return pick(appreciations)
	}
	return ""
}

// MilestoneFor returns a milestone message if streak hits an exact milestone
func MilestoneFor(streak int) string {
	if msg, ok := milestones[streak]; ok {
		return msg
	}
	return ""
}

// PraiseForCommit returns a notification message for any day a commit landed —
// milestone, generic appreciation, or a plain nod for early streak days.
func PraiseForCommit(streak int) string {
	if msg, ok := milestones[streak]; ok {
		return msg
	}
	if streak >= 4 {
		return pick(appreciations)
	}
	return pick(earlyPraise)
}
