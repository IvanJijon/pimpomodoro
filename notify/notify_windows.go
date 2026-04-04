//go:build windows

package notify

import "os/exec"

// Send displays a desktop notification with the given title and message.
func Send(title, message string) {
	script := `[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$textNodes = $template.GetElementsByTagName('text')
$textNodes.Item(0).AppendChild($template.CreateTextNode('` + title + `')) > $null
$textNodes.Item(1).AppendChild($template.CreateTextNode('` + message + `')) > $null
$toast = [Windows.UI.Notifications.ToastNotification]::new($template)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('Pimpomodoro').Show($toast)`
	exec.Command("powershell", "-c", script).Start()
}
