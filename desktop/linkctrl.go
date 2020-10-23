package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"sync"
	"time"
)

type Link struct {
	Cfg      *LinkConfig
	LastFlow  int64
	Bind      string
	Instance *LinkInstance
}

type LinkCtrl struct {
	sync.RWMutex

	Cache   []*Link
}

var linkCtrl *LinkCtrl

func init()  {
	linkCtrl = new(LinkCtrl)
	linkCtrl.Cache = make([]*Link, 0)

	go consoleUpdate()
}

func LinkDelele(binds []string)  {
	linkCtrl.Lock()
	defer linkCtrl.Unlock()

	for _, bind := range binds {
		for i, v := range linkCtrl.Cache {
			if v.Bind != bind {
				continue
			}
			instance := v.Instance
			if instance != nil {
				instance.Close()
			}
			linkCtrl.Cache = append(linkCtrl.Cache[:i], linkCtrl.Cache[i+1:]...)
			break
		}
	}

	syncToFile()
}

func LinkStart(binds []string)  {
	linkCtrl.Lock()
	defer linkCtrl.Unlock()

	for _, bind := range binds {
		for _, v := range linkCtrl.Cache {
			if v.Bind != bind {
				continue
			}
			instance := v.Instance
			if instance == nil {
				instance, err := NewLinkInstance(v.Cfg)
				if err != nil {
					logs.Error(err.Error())
				} else {
					v.Instance = instance
				}
			}
			break
		}
	}
}

func LinkFind(bind string) *LinkConfig {
	linkCtrl.RLock()
	defer linkCtrl.RUnlock()

	for _, v := range linkCtrl.Cache {
		if v.Bind != bind {
			continue
		}
		return v.Cfg
	}
	return nil
}

func LinkStop(binds []string)  {
	linkCtrl.Lock()
	defer linkCtrl.Unlock()

	for _, bind := range binds {
		for _, v := range linkCtrl.Cache {
			if v.Bind != bind {
				continue
			}
			instance := v.Instance
			if instance != nil {
				instance.Close()
				v.Instance = nil
			}
			break
		}
	}
}

func LinkAdd(cfg *LinkConfig) error {
	value, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	logs.Info("link add config : %s", string(value))

	instance, err := NewLinkInstance(cfg)
	if err != nil {
		return err
	}

	bind := fmt.Sprintf("%s:%d", cfg.Iface, cfg.Port)

	linkCtrl.Lock()
	linkCtrl.Cache = append(linkCtrl.Cache, &Link{
		Cfg: cfg, Instance: instance, Bind: bind,
	})
	syncToFile()
	linkCtrl.Unlock()

	return nil
}

func AddLinkItemToConsole(link *Link, idx int) *LinkItem {
	var count int
	var speed int64
	var total int64

	status := STATUS_UNLINK
	if link.Instance != nil {
		count = link.Instance.Channels()
		total = link.Instance.Flows()
		if link.LastFlow < total {
			speed = total - link.LastFlow
		}
		link.LastFlow = total
		status = STATUS_LINK
	}

	return &LinkItem{
		Index: idx,
		Bind: link.Bind,
		Mode: link.Cfg.Mode,
		Count: count,
		Speed: speed/2,
		Size: total,
		Status: status,
	}
}

func consoleUpdate()  {
	time.Sleep(time.Second)
	for  {
		linkCtrl.RLock()
		var output []*LinkItem
		for idx, v := range linkCtrl.Cache {
			output = append(output, AddLinkItemToConsole(v, idx))
		}
		LinkTalbeUpdate(output)
		linkCtrl.RUnlock()
		time.Sleep(2 * time.Second)
	}
}

func syncToFile()  {
	file := fmt.Sprintf("%s\\link.json", appDataDir())

	var output []LinkConfig
	for _, v := range linkCtrl.Cache {
		output = append(output, *v.Cfg)
	}

	value, err := json.Marshal(output)
	if err != nil {
		logs.Error(err.Error())
		return
	}

	err = SaveToFile(file, value)
	if err != nil {
		logs.Error(err.Error())
		return
	}
}

func LinkInit() error {
	file := fmt.Sprintf("%s\\link.json", appDataDir())

	value, err := ioutil.ReadFile(file)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}

	var output []LinkConfig
	err = json.Unmarshal(value, &output)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}

	for _, v := range output {
		temp := v

		instance, err := NewLinkInstance(&temp)
		if err != nil {
			logs.Error(err.Error())
		}

		bind := fmt.Sprintf("%s:%d", temp.Iface, temp.Port)
		linkCtrl.Cache = append(linkCtrl.Cache, &Link{
			Cfg: &temp, Instance: instance, Bind: bind,
		})
	}

	return nil
}

