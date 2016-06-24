package platform

import "syscall"

// ChangeConsoleVisibility Change console visibility
func ChangeConsoleVisibility(visibility bool) {
	var getConsoleWindows = syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	var showWindow = syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	if hwnd, _, _ := showWindow.Call(); hwnd != 0 {
		if visibility {
			sw.Call(hwnd, syscall.SW_RESTORE)
		} else {
			sw.Call(hwnd, syscall.SW_HIDE)
		}
	}
}

// HideConsole Hide console window
func HideConsole(hide bool) {
	ChangeConsoleVisibility(false)
}

// ShowConsole Show console window
func ShowConsole() {
	ChangeConsoleVisibility(true)
}
