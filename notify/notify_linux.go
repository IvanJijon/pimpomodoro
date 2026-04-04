//go:build linux

package notify

import "os/exec"

// Send displays a desktop notification with the given title and message.
func Send(title, message string) {
	if _, err := exec.LookPath("notify-send"); err == nil {
		exec.Command("notify-send", title, message).Start()
	}
}
