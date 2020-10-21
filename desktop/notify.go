package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	"time"
)

var notify *walk.NotifyIcon

func NotifyUpdateFlow(flow string)  {
	if notify == nil {
		return
	}
	notify.SetToolTip(flow)
}

func NotifyExit()  {
	if notify == nil {
		return
	}
	notify.Dispose()
	notify = nil
}

var lastCheck time.Time

func NotifyInit(mw *walk.MainWindow)  {
	NotifyExit()

	var err error

	notify, err = walk.NewNotifyIcon(mw)
	if err != nil {
		logs.Error("new notify icon fail, %s", err.Error())
		return
	}

	err = notify.SetIcon(ICON_Main_Mini)
	if err != nil {
		logs.Error("set notify icon fail, %s", err.Error())
		return
	}

	exitBut := walk.NewAction()
	err = exitBut.SetText("Exit")
	if err != nil {
		logs.Error("notify new action fail, %s", err.Error())
		return
	}

	exitBut.Triggered().Attach(func() {
		MainWindowsClose()
	})

	showBut := walk.NewAction()
	err = showBut.SetText("Show Windows")
	if err != nil {
		logs.Error("notify new action fail, %s", err.Error())
		return
	}

	showBut.Triggered().Attach(func() {
		MainWindowsVisible(true)
	})

	if err := notify.ContextMenu().Actions().Add(showBut); err != nil {
		logs.Error("notify add action fail, %s", err.Error())
		return
	}

	if err := notify.ContextMenu().Actions().Add(exitBut); err != nil {
		logs.Error("notify add action fail, %s", err.Error())
		return
	}

	notify.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		now := time.Now()
		if now.Sub(lastCheck) < 2 * time.Second {
			MainWindowsVisible(true)
		}
		lastCheck = now
	})

	notify.SetVisible(true)
}

