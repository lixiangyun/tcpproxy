package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func ShowToolBar(cfg * LinkConfig)  {
	var dlg *walk.Dialog
	var acceptPB *walk.PushButton
	var backendView *walk.TableView

	backendTable := new(BackendModel)
	backendTable.Input(cfg.Backend)

	cnt, err := Dialog{
		AssignTo: &dlg,
		Title: "Link Detail",
		Icon: walk.IconInformation(),
		DefaultButton: &acceptPB,
		Size: Size{400, 300},
		MinSize: Size{400, 300},
		Layout:  VBox{
			Alignment: AlignHNearVCenter,
			Margins: Margins{Top: 10, Bottom: 10, Left: 10, Right: 10},
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Bind Address:",
					},
					Label{
						Text: fmt.Sprintf("%s:%d",
							cfg.Iface, cfg.Port),
					},
					Label{
						Text: "Bind Timeout:",
					},
					Label{
						Text: fmt.Sprintf("%d Second", cfg.Timeout),
					},
					Label{
						Text: "Load Balance:",
					},
					Label{
						Text: cfg.Mode,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					Label{
						Text: "Backend List:",
					},
					TableView{
						AssignTo: &backendView,
						AlternatingRowBG: true,
						ColumnsOrderable: true,
						Columns: []TableViewColumn{
							{Title: "#", Width: 20},
							{Title: "Address", Width: 110},
							{Title: "Timeout", Width: 50},
							{Title: "Weight", Width: 50},
							{Title: "Main/Standby", Width: 80},
						},
						StyleCell: func(style *walk.CellStyle) {
							if style.Row()%2 == 0 {
								style.BackgroundColor = walk.RGB(248, 248, 255)
							} else {
								style.BackgroundColor = walk.RGB(220, 220, 220)
							}
						},
						Model: backendTable,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text: "OK",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(MainWindowsCtrl())
	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("show link dialog return %d", cnt)
	}
}
