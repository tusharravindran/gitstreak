# gitstreak

Keep your GitHub contribution streak alive. Get your streak stats in the terminal and a macOS notification at 9 PM if you haven't committed yet.

```
$ gitstreak status --username tusharravindran

  Tushar Ravindran
  github.com/tusharravindran

  ✅ You've committed today  (2 contributions)

  Current Streak       Longest Streak       This Year
  🔥 2 days            ⚡ 14 days           📊 63 commits

  Last 4 weeks  ░░░░░▒▒▓▒░░▒▒▒▓▓▒░░░▒▒▓▒░░▒▒
```

## Install

```bash
brew tap tusharravindran/homebrew-devdoctor
brew install gitstreak
```

Or download a binary from [Releases](https://github.com/tusharravindran/gitstreak/releases).

## Setup

You need a GitHub personal access token with `read:user` scope.

1. Create one at https://github.com/settings/tokens
2. Export it:

```bash
export GITHUB_TOKEN=your_token_here
export GITHUB_USERNAME=your_username
```

Add both to your `~/.zshrc` or `~/.bashrc` to make them permanent.

## Usage

### Check your streak

```bash
gitstreak status --username <your-github-handle>
```

### Get task suggestions

No commit ideas? Run:

```bash
gitstreak suggest
```

### Enable 9 PM reminder

Install a background reminder that sends a macOS notification at 9 PM if you haven't committed today:

```bash
gitstreak watch --username <your-github-handle>
```

Remove it:

```bash
gitstreak unwatch
```

### Commands

| Command | Description |
|---|---|
| `gitstreak status` | Show streak, heatmap, and task suggestions |
| `gitstreak suggest` | List 5 quick tasks to commit today |
| `gitstreak watch` | Install 9 PM daily macOS notification |
| `gitstreak unwatch` | Remove the daily reminder |

## Environment variables

| Variable | Description |
|---|---|
| `GITHUB_TOKEN` | GitHub personal access token (`read:user` scope) |
| `GITHUB_USERNAME` | Your GitHub handle (used as default when `--username` is omitted) |

## Why

I kept breaking my GitHub streak because I forgot to commit on weekends or during deep work days. A simple terminal check + a phone notification fixed it. Built this in a weekend, figured others would want it too.

## License

MIT
