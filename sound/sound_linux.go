//go:build linux

package sound

import (
	"fmt"
	"os"
	"os/exec"
)

// PlayAlarm plays a system sound using the first available player and sound file.
func PlayAlarm() {
	files := []string{
		"/usr/share/sounds/freedesktop/stereo/complete.oga",
		"/usr/share/sounds/freedesktop/stereo/bell.oga",
	}

	players := []string{"paplay", "aplay"}

	for _, player := range players {
		if _, err := exec.LookPath(player); err != nil {
			continue
		}
		for _, file := range files {
			if _, err := os.Stat(file); err == nil {
				exec.Command(player, file).Start()
				return
			}
		}
	}

	fmt.Print("\a")
}
