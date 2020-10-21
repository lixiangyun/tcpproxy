package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func boxAction(from walk.Form, title string, icon *walk.Icon, message string)  {
	var dlg *walk.Dialog
	var cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo: &dlg,
		Title: title,
		Icon: icon,
		MinSize: Size{210, 150},
		Size: Size{210, 150},
		MaxSize: Size{310, 210},
		CancelButton: &cancelPB,
		Layout:  VBox{},
		Children: []Widget{
			TextLabel{
				Text: message,
				MinSize: Size{200, 150},
				MaxSize: Size{300, 200},
			},
			PushButton{
				AssignTo:  &cancelPB,
				Text:      "OK",
				OnClicked: func() {
					dlg.Accept()
				},
			},
		},
	}.Run(from)

	if err != nil {
		logs.Error(err.Error())
	}
}

func ErrorBoxAction(form walk.Form, message string) {
	boxAction(form, "Error", walk.IconError(), message)
}

func InfoBoxAction(form walk.Form, message string) {
	boxAction(form, "Info", walk.IconInformation(), message)
}

func ConfirmBoxAction(form walk.Form, message string) {
	boxAction(form, "Warning", walk.IconWarning(), message)
}
