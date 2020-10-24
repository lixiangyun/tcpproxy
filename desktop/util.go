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
	"strconv"
	"strings"
	"syscall"
	"time"

	mathrand "math/rand"
)

func VersionGet() string {
	return "v1.0.0"
}

func ListenCheck(addr string, port int) bool {
	list, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	defer list.Close()
	return true
}

func SaveToFile(name string, body []byte) error {
	return ioutil.WriteFile(name, body, 0664)
}

func InterfaceAddsGet(iface *net.Interface) ([]net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, nil
	}
	ips := make([]net.IP, 0)
	for _, v:= range addrs {
		ipone, _, err:= net.ParseCIDR(v.String())
		if err != nil {
			continue
		}
		if len(ipone) > 0 {
			ips = append(ips, ipone)
		}
	}
	return ips, nil
}

func AddressValid(addr string) bool {
	if addr == "" {
		return false
	}
	list := strings.Split(addr,":")
	if len(list) != 2 {
		logs.Error("address valid fail, %s", addr)
		return false
	}
	ip := net.ParseIP(list[0])
	if ip == nil {
		logs.Error("address valid fail, %s", addr)
		return false
	}
	cnt, err := strconv.Atoi(list[1])
	if err != nil {
		logs.Error("address valid fail, %s", err.Error())
		return false
	}
	if cnt > 65535 || cnt < 0 {
		logs.Error("address valid fail, %s", addr)
		return false
	}
	return true
}

func IsIPv4(ip net.IP) bool {
	return strings.Index(ip.String(), ".") != -1
}

func InterfaceLocalIP(inface *net.Interface) ([]net.IP, error) {
	addrs, err := InterfaceAddsGet(inface)
	if err != nil {
		return nil, err
	}
	var output []net.IP
	for _, v := range addrs {
		if IsIPv4(v) == true {
			output = append(output, v)
		}
	}
	if len(output) == 0 {
		return nil, fmt.Errorf("interface not ipv4 address.")
	}
	return output, nil
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

func init()  {
	mathrand.Seed(time.Now().Unix())
}