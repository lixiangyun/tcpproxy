package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	"os"
)

func IconLoadFromBox(filename string, size walk.Size) *walk.Icon {
	body, err := BoxFile().Bytes(filename)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	dir := DEFAULT_HOME + "\\icon\\"
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 644)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}
	}
	filepath := dir + filename
	err = SaveToFile(filepath, body)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	icon, err := walk.NewIconFromFileWithSize(filepath, size)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	return icon
}


func IconLoadImageFromBox(filename string) walk.Image {
	body, err := BoxFile().Bytes(filename)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	dir := DEFAULT_HOME + "\\image\\"
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 644)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}
	}
	filepath := dir + filename
	err = SaveToFile(filepath, body)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	image, err := walk.NewImageFromFile(filepath)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	return image
}

var ICON_Main          *walk.Icon
var ICON_Main_Mini     *walk.Icon

var ICON_TOOL_ADD      *walk.Icon
var ICON_TOOL_DEL      *walk.Icon
var ICON_TOOL_LINK     *walk.Icon
var ICON_TOOL_UNLINK   *walk.Icon
var ICON_TOOL_SETTING  *walk.Icon

var ICON_STATUS_UNLINK *walk.Icon
var ICON_STATUS_LINK *walk.Icon

var ICON_Max_Size = walk.Size{
	Width: 128, Height: 128,
}

var ICON_Tool_Size = walk.Size{
	Width: 64, Height: 64,
}

var ICON_Min_Size = walk.Size{
	Width: 24, Height: 24,
}

func IconInit() error {
	ICON_Main = IconLoadFromBox("main.ico", ICON_Max_Size)
	ICON_Main_Mini = IconLoadFromBox("mainmini.ico", ICON_Min_Size)

	ICON_TOOL_ADD = IconLoadFromBox("add.ico", ICON_Tool_Size)
	ICON_TOOL_DEL = IconLoadFromBox("delete.ico", ICON_Tool_Size)
	ICON_TOOL_LINK = IconLoadFromBox("link.ico", ICON_Tool_Size)
	ICON_TOOL_UNLINK = IconLoadFromBox("unlink.ico", ICON_Tool_Size)
	ICON_TOOL_SETTING = IconLoadFromBox("setting.ico", ICON_Tool_Size)

	ICON_STATUS_UNLINK = IconLoadFromBox("status_unlink.ico", ICON_Min_Size)
	ICON_STATUS_LINK = IconLoadFromBox("status_link.ico", ICON_Min_Size)

	return nil
}

