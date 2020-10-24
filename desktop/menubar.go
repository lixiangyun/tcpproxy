package main

import (
	. "github.com/lxn/walk/declarative"
)

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Action{
			Text: "Runlog",
			OnTriggered: func() {
				OpenBrowserWeb(LogDirGet())
			},
		},
		Action{
			Text: "Mini Windows",
			OnTriggered: func() {
				MainWindowsVisible(false)
			},
		},
		Action{
			Text: "About",
			OnTriggered: func() {
				AboutAction(mainWindowCtrl.ctrl)
			},
		},
	}
}
