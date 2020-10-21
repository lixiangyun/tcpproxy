package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	mathrand "math/rand"
)

func VersionGet() string {
	return "v1.0.0"
}

func IsConnect(address string, timeout int) bool {
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func GetTimeStamp() string {
	now := time.Now()
	year, month, day := now.Date()
	return fmt.Sprintf(
		"%4d-%02d-%02d %02d:%02d:%02d",
		year, month, day, now.Hour(), now.Minute(), now.Second())
}

func GetTimeStampNumber() string {
	now := time.Now()
	year, month, day := now.Date()
	return fmt.Sprintf(
		"%4d%02d%02d%02d%02d%02d.%03d",
		year, month, day,
		now.Hour(), now.Minute(), now.Second(),
		time.Duration(now.Nanosecond())/time.Millisecond)
}

func SaveToFile(name string, body []byte) error {
	return ioutil.WriteFile(name, body, 0664)
}

func GetToken(length int) string {
	token := make([]byte, length)
	bytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!#$%^&*"
	for i:=0; i<length; i++  {
		token[i] = bytes[mathrand.Int()%len(bytes)]
	}
	return string(token)
}

func GetUser(length int) string {
	token := make([]byte, length)
	bytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i:=0; i<length; i++  {
		token[i] = bytes[mathrand.Int()%len(bytes)]
	}
	return string(token)
}

func CapSignal(proc func())  {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	go func() {
		sig := <- signalChan
		proc()
		logs.Error("recv signcal %s, ready to exit", sig.String())
		os.Exit(-1)
	}()
}

func ByteViewLite(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%db", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.1fKb", float64(size)/float64(1024))
	} else {
		return fmt.Sprintf("%.1fMb", float64(size)/float64(1024*1024))
	}
}

func ByteView(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.1fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fGB", float64(size)/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.1fTB", float64(size)/float64(1024*1024*1024*1024))
	}
}

func StringList(list []string) string {
	var body string
	for idx,v := range list {
		if idx == len(list) - 1 {
			body += fmt.Sprintf("%s",v)
		}else {
			body += fmt.Sprintf("%s;",v)
		}
	}
	return body
}

type logconfig struct {
	Filename string  `json:"filename"`
	Level    int     `json:"level"`
	MaxLines int     `json:"maxlines"`
	MaxSize  int     `json:"maxsize"`
	Daily    bool    `json:"daily"`
	MaxDays  int     `json:"maxdays"`
	Color    bool    `json:"color"`
}

var logCfg = logconfig{Filename: os.Args[0], Level: 7, Daily: true, MaxDays: 30, Color: true}

func LogInit() error {
	logCfg.Filename = fmt.Sprintf("%s%c%s", LogDirGet(), os.PathSeparator, "runlog.log")
	value, err := json.Marshal(&logCfg)
	if err != nil {
		return err
	}

	if DebugFlag() {
		err = logs.SetLogger(logs.AdapterConsole)
	} else {
		err = logs.SetLogger(logs.AdapterFile, string(value))
	}

	if err != nil {
		return err
	}

	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	return nil
}

func StringDiff(oldlist []string, newlist []string) ([]string, []string) {
	del := make([]string, 0)
	add := make([]string, 0)
	for _,v1 := range oldlist {
		flag := false
		for _,v2 := range newlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if flag == false {
			del = append(del, v1)
		}
	}
	for _,v1 := range newlist {
		flag := false
		for _,v2 := range oldlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if flag == false {
			add = append(add, v1)
		}
	}
	return del, add
}

func StringClone(list []string) []string {
	output := make([]string, len(list))
	copy(output, list)
	return output
}

func WriteFull(w io.Writer, body []byte) error {
	begin := 0
	for  {
		cnt, err := w.Write(body[begin:])
		if cnt > 0 {
			begin += cnt
		}
		if begin >= len(body) {
			return err
		}
		if err != nil {
			return err
		}
	}
}

func TouchDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 644)
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		return nil
	}
	if !info.IsDir() {
		return fmt.Errorf("[%s] is not directory", dir)
	}
	return nil
}

func init()  {
	mathrand.Seed(time.Now().Unix())
}