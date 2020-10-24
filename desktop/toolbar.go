package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var toolBars *walk.ToolBar

func LinkDelToolBar()  {
	list := consoleLinkTable.LinkTableSelectList()
	if len(list) == 0 {
		ErrorBoxAction(MainWindowsCtrl(), "No object selected")
		return
	}
	LinkDelele(list)
	consoleLinkTable.LinkTableSelectClean()
}

func LinkStartToolBar()  {
	list := consoleLinkTable.LinkTableSelectList()
	if len(list) == 0 {
		ErrorBoxAction(MainWindowsCtrl(), "No object selected")
		return
	}
	LinkStart(list)
	consoleLinkTable.LinkTableSelectClean()
}

func LinkStopToolBar()  {
	list := consoleLinkTable.LinkTableSelectList()
	if len(list) == 0 {
		ErrorBoxAction(MainWindowsCtrl(), "No object selected")
		return
	}
	LinkStop(list)
	consoleLinkTable.LinkTableSelectClean()
}

func ToolBarInit() ToolBar {
	return ToolBar{
		AssignTo: &toolBars,
		ButtonStyle: ToolBarButtonImageOnly,
		MinSize: Size{Width: 64, Height: 64},
		Items: []MenuItem{
			Action{
				Text: "Add Link",
				Image: ICON_TOOL_ADD,
				OnTriggered: func() {
					AddToolBar()
				},
			},
			Action{
				Text: "Delete Link",
				Image: ICON_TOOL_DEL,
				OnTriggered: func() {
					go LinkDelToolBar()
				},
			},
			Action{
				Text: "Link",
				Image: ICON_TOOL_LINK,
				OnTriggered: func() {
					go LinkStartToolBar()
				},
			},
			Action{
				Text: "Unlink",
				Image: ICON_TOOL_UNLINK,
				OnTriggered: func() {
					go LinkStopToolBar()
				},
			},
			//Action{
			//	Text: "Setting",
			//	Image: ICON_TOOL_SETTING,
			//	OnTriggered: func() {
			//		//Setting()
			//	},
			//},
		},
	}
}
