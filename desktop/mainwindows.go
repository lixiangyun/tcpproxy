package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
	"time"
)


var mainWindowWidth = 700
var mainWindowHeight = 400

type MainWindowCtrl struct {
	instance *MainWindow
	ctrl     *walk.MainWindow
	exitInt  int
	exit     chan struct{}
}

var mainWindowCtrl *MainWindowCtrl

func init() {
	mainWindowCtrl = new(MainWindowCtrl)
	mainWindowCtrl.exit = make(chan struct{}, 10)
}

func MainWindowsVisible(flag bool)  {
	if mainWindowCtrl.ctrl != nil {
		mainWindowCtrl.ctrl.SetVisible(flag)
	}
}

func MainWindowsClose()  {
	logs.Info("main windows ready to close")
	mainWindowCtrl.exitInt = 0
	mainWindowCtrl.exit <- struct{}{}
}

func AppExitPreFunc()  {
	if mainWindowCtrl.ctrl != nil {
		mainWindowCtrl.ctrl.Close()
		mainWindowCtrl.ctrl = nil
	}
	NotifyExit()
	if err:= recover();err != nil{
		logs.Error(err)
	}
}

func MainWindowsCtrl() *walk.MainWindow {
	return mainWindowCtrl.ctrl
}

func MainWindowsExit()  {
	CapSignal(AppExitPreFunc)
	<- mainWindowCtrl.exit
	AppExitPreFunc()
	os.Exit(mainWindowCtrl.exitInt)
}

func MainWindowStart() error {
	logs.Info("main windows start")
	mainWindowCtrl.instance = mainWindowBuilder(&mainWindowCtrl.ctrl)

	go func() {
		cnt, err := mainWindowCtrl.instance.Run()
		if err != nil {
			logs.Error(err.Error())
		}
		mainWindowCtrl.exitInt = cnt
		mainWindowCtrl.exit <- struct{}{}
		logs.Info("main windows close")
	}()

	for  {
		if mainWindowCtrl.ctrl != nil && mainWindowCtrl.ctrl.Visible() {
			mainWindowCtrl.ctrl.SetSize(walk.Size{
				mainWindowWidth,
				mainWindowHeight})
			break
		}
		time.Sleep(200*time.Millisecond)
	}

	NotifyInit(mainWindowCtrl.ctrl)

	logs.Info("main windows start success")
	return nil
}

func mainWindowBuilder(mw **walk.MainWindow) *MainWindow {
	return &MainWindow{
		Title:   "TcpProxy " + VersionGet(),
		Icon: ICON_Main,
		AssignTo: mw,
		MinSize: Size{mainWindowWidth, mainWindowHeight-1},
		Size: Size{mainWindowWidth, mainWindowHeight-1},
		Layout:  VBox{
			Alignment: AlignHNearVNear,
			MarginsZero: true,
			Margins: Margins{Left: 10, Top: 10},
		},
		MenuItems: MenuBarInit(),
		Children: []Widget{
			Composite{
				Layout: HBox{ MarginsZero: true},
				Children: []Widget{
					ToolBarInit(),
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, Margins: Margins{Right: 10, Bottom: 10}},
				Children: TableWight(),
			},
		},
	}
}
