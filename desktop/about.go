package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os/exec"
)

func OpenBrowserWeb(url string)  {
	cmd := exec.Command("rundll32","url.dll,FileProtocolHandler", url)
	err := cmd.Run()
	if err != nil {
		logs.Error("run cmd fail, %s", err.Error())
	}
}

var aboutContext = ""

var image1 walk.Image
var image2 walk.Image

func AboutAction( mw *walk.MainWindow ) {
	var ok    *walk.PushButton
	var about *walk.Dialog
	var err error

	if aboutContext == "" {
		temp, err := BoxFile().String("about.txt")
		if err != nil {
			logs.Error(err.Error())
		}
		aboutContext = temp
	}

	if image1 == nil {
		image1 = IconLoadImageFromBox("sponsor1.jpg")
	}

	if image2 == nil {
		image2 = IconLoadImageFromBox("sponsor2.jpg")
	}

	_, err = Dialog{
		AssignTo:      &about,
		Title:         "About",
		Icon:          walk.IconInformation(),
		MinSize:       Size{Width: 300, Height: 200},
		DefaultButton: &ok,
		Layout:  VBox{},
		Children: []Widget{
			TextLabel{
				Text: aboutContext,
				MinSize:       Size{Width: 250, Height: 200},
				MaxSize:       Size{Width: 290, Height: 400},
			},
			Label{
				Text: "Version: "+ VersionGet(),
				TextAlignment: AlignCenter,
			},
			VSpacer{
				MinSize: Size{Height: 10},
			},
			Label{
				Text: "Sponsor",
				TextAlignment: AlignCenter,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{
						MinSize: Size{Width: 10},
					},
					ImageView{
						ToolTipText: "Ali Pay",
						Image:    image1,
						MaxSize:  Size{80, 80},
					},
					HSpacer{
						MinSize: Size{Width: 10},
					},
					ImageView{
						ToolTipText: "Wecart Pay",
						Image:    image2,
						MaxSize:  Size{80, 80},
					},
					HSpacer{
						MinSize: Size{Width: 10},
					},
				},
			},
			PushButton{
				Text:      "Paypal.me",
				OnClicked: func() {
					OpenBrowserWeb("https://paypal.me/lixiangyun")
				},
			},
			PushButton{
				Text:      "Official Web",
				OnClicked: func() {
					OpenBrowserWeb("https://github.com/lixiangyun/tcpproxy")
				},
			},
			PushButton{
				Text:      "OK",
				OnClicked: func() { about.Cancel() },
			},
		},
	}.Run(mw)

	if err != nil {
		logs.Error(err.Error())
	}
}
