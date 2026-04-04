//go:build darwin

package sound

import "os/exec"

// PlayAlarm plays a system sound to notify the user.
func PlayAlarm() {
	exec.Command("afplay", "/System/Library/Sounds/Ping.aiff").Start()
}
