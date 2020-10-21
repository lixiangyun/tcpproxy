package main

import "github.com/astaxie/beego/logs"

func main()  {
	err := DebugInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = FileInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = LogInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = BoxInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = IconInit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	err = MainWindowStart()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	MainWindowsExit()
}