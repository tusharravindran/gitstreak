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

I kept breaking my GitHub streak because I forgot to commit on weekends or during deep-focus days at work. A quick terminal check and a 9 PM notification fixed the forgetting. Skip days fixed the guilt. Built this in a weekend.

---

## License

MIT — [Tushar Ravindran](https://tusharravindran.github.io)
