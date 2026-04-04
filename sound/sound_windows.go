//go:build windows

package sound

import "os/exec"

// PlayAlarm plays a system sound to notify the user.
func PlayAlarm() {
	exec.Command("powershell", "-c", "[System.Media.SystemSounds]::Beep.Play()").Start()
}
