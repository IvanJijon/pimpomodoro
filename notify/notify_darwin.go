//go:build darwin

package notify

import "os/exec"

// Send displays a desktop notification with the given title and message.
func Send(title, message string) {
	exec.Command("osascript", "-e", `display notification "`+message+`" with title "`+title+`"`).Start()
}
