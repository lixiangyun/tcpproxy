package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var toolBars *walk.ToolBar

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
					//AddJobOnce()
				},
			},
			Action{
				Text: "Delete Link",
				Image: ICON_TOOL_DEL,
				OnTriggered: func() {
					//AddJobBatch()
				},
			},
			Action{
				Text: "Link",
				Image: ICON_TOOL_LINK,
				OnTriggered: func() {
					//KeepSet()
				},
			},
			Action{
				Text: "Unlink",
				Image: ICON_TOOL_DEL,
				OnTriggered: func() {

				},
			},
			Action{
				Text: "Setting",
				Image: ICON_TOOL_SETTING,
				OnTriggered: func() {
					//Setting()
				},
			},
		},
	}
}
