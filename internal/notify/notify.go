package notify

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.tusharravindran.gitstreak</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.BinaryPath}}</string>
        <string>notify-check</string>
        <string>--username</string>
        <string>{{.Username}}</string>
    </array>
    <key>StartCalendarInterval</key>
    <dict>
        <key>Hour</key>
        <integer>{{.Hour}}</integer>
        <key>Minute</key>
        <integer>{{.Minute}}</integer>
    </dict>
    <key>EnvironmentVariables</key>
    <dict>
        <key>GITSTREAK_GH_TOKEN</key>
        <string>{{.Token}}</string>
    </dict>
    <key>RunAtLoad</key>
    <false/>
    <key>StandardOutPath</key>
    <string>{{.LogPath}}</string>
    <key>StandardErrorPath</key>
    <string>{{.LogPath}}</string>
</dict>
</plist>
`

type PlistData struct {
	BinaryPath string
	Username   string
	Token      string
	LogPath    string
	Hour       int
	Minute     int
}

func Install(username string, hour, minute int) error {
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find binary path: %w", err)
	}

	token := os.Getenv("GITSTREAK_GH_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}

	home, _ := os.UserHomeDir()
	plistPath := filepath.Join(home, "Library", "LaunchAgents", "com.tusharravindran.gitstreak.plist")
	logPath := filepath.Join(home, "Library", "Logs", "gitstreak.log")

	data := PlistData{
		BinaryPath: binaryPath,
		Username:   username,
		Token:      token,
		LogPath:    logPath,
		Hour:       hour,
		Minute:     minute,
	}

	f, err := os.Create(plistPath)
	if err != nil {
		return fmt.Errorf("could not write plist: %w", err)
	}
	defer f.Close()

	tmpl := template.Must(template.New("plist").Parse(plistTemplate))
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("could not render plist: %w", err)
	}

	// The plist embeds the GitHub token in plaintext — restrict it to the
	// owner so other local users/processes can't read it off disk.
	if err := os.Chmod(plistPath, 0600); err != nil {
		return fmt.Errorf("could not set plist permissions: %w", err)
	}

	// Unload any previously-running job first — launchctl load is a no-op on an
	// already-loaded job, so a stale schedule would otherwise stick around even
	// after the plist file itself is rewritten.
	exec.Command("launchctl", "unload", plistPath).Run()

	cmd := exec.Command("launchctl", "load", "-w", plistPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl load failed: %s", string(out))
	}

	fmt.Printf("✅ gitstreak reminder installed\n")
	fmt.Printf("   Plist: %s\n", plistPath)
	fmt.Printf("   Logs:  %s\n", logPath)
	return nil
}

// IsInstalled reports whether the daily reminder LaunchAgent is currently set up.
func IsInstalled() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	plistPath := filepath.Join(home, "Library", "LaunchAgents", "com.tusharravindran.gitstreak.plist")
	_, err = os.Stat(plistPath)
	return err == nil
}

func Uninstall() error {
	home, _ := os.UserHomeDir()
	plistPath := filepath.Join(home, "Library", "LaunchAgents", "com.tusharravindran.gitstreak.plist")

	cmd := exec.Command("launchctl", "unload", plistPath)
	cmd.Run()

	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not remove plist: %w", err)
	}

	fmt.Println("✅ Daily reminder uninstalled")
	return nil
}

func Send(title, message string) {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "Funk"`, message, title)
	exec.Command("osascript", "-e", script).Run()
}

// SendUrgent fires a louder notification + voice for roast/critical alerts
func SendUrgent(title, message, voice string) {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "Sosumi"`, message, title)
	exec.Command("osascript", "-e", script).Run()
	speak(message, voice)
}

// SendWithVoice fires the standard pleasant notification and speaks the
// message aloud, using the same voice as SendUrgent for consistency.
func SendWithVoice(title, message, voice string) {
	Send(title, message)
	speak(message, voice)
}

func speak(message, voice string) {
	if voice == "" {
		voice = "Samantha"
	}
	exec.Command("say", "-v", voice, message).Run()
}
