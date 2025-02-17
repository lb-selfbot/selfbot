package interop

// I/O

func ClearConsole() {
	logging.Console.Clear()
	logging.UI.Clear()

	logging.Logging.PrintBanner(true, true)
}

func Notification(message string) {
	logging.UI.Notification(message)
}
