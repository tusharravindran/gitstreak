# gitstreak

A CLI tool to keep your GitHub contribution streak alive.

See your streak in the terminal, get quick commit ideas when you're stuck, and set a daily reminder before you forget. Skip weekends without breaking your streak.

```
$ gitstreak status --username tusharravindran

  Tushar Ravindran
  github.com/tusharravindran

  ✅ You've committed today  (3 contributions)

  Current Streak       Longest Streak       This Year
  🔥 14 days           ⚡ 30 days           📊 180 commits

  Last 4 weeks  ░▒▒▓▒░░▒▒▒▓▓▒░░░▒▒▓▒░░▒▒▓▒░░
                ░ 0   ▒ 1–2   ▓ 3–5   █ 6+

  🔥 14 days and still going. The momentum is real.
```

---

## Install

```bash
brew tap tusharravindran/homebrew-devdoctor
brew install gitstreak
```

Or grab a binary from [Releases](https://github.com/tusharravindran/gitstreak/releases).

---

## Setup

You need a GitHub personal access token with `read:user` scope.

1. Create one at https://github.com/settings/tokens (fine-grained, `read:user` only)
2. Add to your shell config (`~/.zshrc` or `~/.bashrc`):

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
export GITHUB_USERNAME=your_github_handle
```

3. Reload: `source ~/.zshrc`

`read:user` is all you need for `status`, `suggest`, and `watch`. If you also want
cheat-day detection (`gitstreak audit`, and the inline nudge in `status`), the token
needs `repo` scope instead — that's what lets gitstreak read commit contents in your
private repos. Public-only accounts can skip this; it's opt-in.

---

## Commands

### `gitstreak status`

Shows your current streak, heatmap, and suggestions if you haven't committed today.

```bash
gitstreak status
gitstreak status --username tusharravindran
```

Streak appreciation fires at 4+ days. Milestone messages at 7, 14, 30, 50, 100, 365 days.
If your streak has been broken for more than 3 days, it roasts you instead.

If today's commit looks like it was just there to keep the streak alive — a one-line
change, a lazy commit message, docs/config-only, or the same single file edited days in
a row — `status` roasts that too, instead of showing the usual streak appreciation.
Requires a `repo`-scoped token (see Setup); silently skipped otherwise.

---

### `gitstreak audit`

Checks whether your recent streak days were genuine or cheat days — a one-line commit
that technically kept the streak alive but didn't really earn it.

```bash
gitstreak audit
gitstreak audit --days 14
```

```
  Auditing last 7 days for @tusharravindran

  2026-06-30  ✅ genuine
  2026-07-01  😏 cheat day  (1 line changed, message: "wip")
  2026-07-02  ✅ genuine
  2026-07-03  —  no commits
  2026-07-04  😏 cheat day  (docs/config only, no source changes)
  2026-07-05  ✅ genuine
  2026-07-06  ✅ genuine

  4 genuine, 2 cheat days, 0 unauditable
  Your real streak is a bit shorter than your GitHub graph says 👀
```

Days it can't inspect (private repos the token can't read, or GitHub-anonymized
contributions) show as unauditable rather than being guessed at either way.

Requires a `repo`-scoped `GITHUB_TOKEN` (see Setup) — `read:user` alone can't read
commit contents.

---

### `gitstreak suggest`

Get 5 quick task ideas you can commit today — even on low-energy days.

```bash
gitstreak suggest
```

```
  💡 Things you can commit today:

  1.  📝  Write a README
       Add or improve the README for one of your repos. Counts as a commit.
       ⏱ 10 min

  2.  🧪  Write a test
       Pick an untested function and write one test. One is better than zero.
       ⏱ 15 min
```

---

### `gitstreak watch`

Installs a background reminder using macOS launchd. Checks every day at your configured time and sends a notification if you haven't committed.

```bash
gitstreak watch
gitstreak watch --username tusharravindran
```

```bash
gitstreak unwatch   # remove the reminder
```

The reminder is skipped automatically on your configured skip days.

---

### `gitstreak config`

View or update your settings.

```bash
gitstreak config                              # view current settings
gitstreak config --username tusharravindran   # set default username
gitstreak config --reminder-time 20:00        # set reminder to 8 PM
gitstreak config --reminder-time 8:30pm       # 12-hour format also works
gitstreak config --skip-days Sat,Sun          # skip weekends
gitstreak config --skip-days Sat,Sun,Fri      # skip Fri–Sun
gitstreak config --clear-skip-days            # go back to daily reminders
```

After changing reminder time or skip days, re-run `gitstreak watch` to apply.

**Skip days explained:**
- On a skip day, no reminder fires at your configured time
- If you commit on a skip day anyway, it counts toward your streak normally
- If you don't commit on a skip day, it does **not** break your streak
- Skipping more than 3 consecutive days will earn you a roast 😬

---

## Settings

Settings are stored at `~/.config/gitstreak/config.json`.

| Setting | Default | Description |
|---|---|---|
| `username` | — | Default GitHub handle (avoids typing `--username` every time) |
| `reminder_hour` | `21` | Hour for daily reminder (24h) |
| `reminder_minute` | `0` | Minute for daily reminder |
| `skip_days` | `[]` | Weekdays to skip reminders and exclude from streak (0=Sun, 6=Sat) |

---

## Environment variables

| Variable | Description |
|---|---|
| `GITHUB_TOKEN` | GitHub PAT with `read:user` scope |
| `GITHUB_USERNAME` | Default username (overridden by `gitstreak config --username`) |

---

## Why

For a long time I was heads-down on my organization's GitHub — shipping features, closing PRs, doing real work. But all of that lives on a company account. My personal profile was a ghost town.

I decided to block 2–3 hours every day specifically for my personal GitHub — learning, building, and shipping tools that are actually mine. The problem was I kept forgetting. End of the day rolls around, I've been coding for 8 hours, and the last thing on my mind is opening a side project.

So I built gitstreak. A quick `gitstreak status` in the morning tells me where I stand. A nudge at a time I choose reminds me before the day slips away. Skip days mean I'm not guilted on weekends. And when the streak grows, it actually feels like something worth protecting.

Built this in a weekend. Using it every day.

---

## License

MIT — [Tushar Ravindran](https://tusharravindran.github.io)
