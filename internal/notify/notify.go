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
        <integer>21</integer>
        <key>Minute</key>
        <integer>0</integer>
    </dict>
    <key>EnvironmentVariables</key>
    <dict>
        <key>GITHUB_TOKEN</key>
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
}

func Install(username string) error {
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find binary path: %w", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
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

	cmd := exec.Command("launchctl", "load", "-w", plistPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl load failed: %s", string(out))
	}

	fmt.Printf("✅ gitstreak will remind you at 9:00 PM daily\n")
	fmt.Printf("   Plist: %s\n", plistPath)
	fmt.Printf("   Logs:  %s\n", logPath)
	return nil
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
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "Glass"`, message, title)
	exec.Command("osascript", "-e", script).Run()
}
