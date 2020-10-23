package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type LinkChannel struct {
	key    string
	remote net.Conn
	proxy  net.Conn
	resvflow int64
	sendflow int64
}

type LinkInstance struct {
	sync.RWMutex
	sync.WaitGroup

	close bool

	addr string
	lb   LoadBalance
	cfg *LinkConfig
	list net.Listener
	channels map[string]*LinkChannel
}

func NewLinkInstance(item *LinkConfig) (*LinkInstance, error) {
	address := fmt.Sprintf("%s:%d", item.Iface, item.Port)
	list, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	link := new(LinkInstance)
	link.addr = address
	link.list = list
	link.channels = make(map[string]*LinkChannel, 1024)
	link.cfg = item
	link.lb = NewLoadBalance(item.Mode, item.Backend)

	link.Add(1)
	go link.start()
	return link, nil
}

func (l *LinkInstance)proxy(wg *sync.WaitGroup, conn1 net.Conn)  {
	var err error
	var conn2 net.Conn

	defer func() {
		wg.Done()
		conn1.Close()
		if conn2 != nil {
			conn2.Close()
		}
	}()

	key   := conn1.RemoteAddr().String()
	idx   := l.lb.Next(key)
	proxy := l.cfg.Backend[idx]

	if proxy.Timeout == 0 {
		conn2, err = net.Dial("tcp", proxy.Address )
	} else {
		timeout := time.Second * time.Duration(proxy.Timeout)
		conn2, err = net.DialTimeout("tcp", proxy.Address, timeout)
	}

	if err != nil {
		logs.Error(err.Error())
		return
	}

	channel := new(LinkChannel)
	channel.remote = conn1
	channel.proxy = conn2
	channel.key = key

	l.Lock()
	l.channels[key] = channel
	l.Unlock()

	wg2 := new(sync.WaitGroup)
	wg2.Add(2)
	go connect(wg2, conn1, conn2, &channel.sendflow)
	go connect(wg2, conn1, conn2, &channel.resvflow)
	wg2.Wait()

	l.Lock()
	delete(l.channels, key)
	l.Unlock()
}

func (l *LinkInstance)start()  {
	defer l.Done()
	logs.Info("link instance %s start", l.addr)

	wg := new(sync.WaitGroup)
	for  {
		if l.close {
			break
		}
		conn, err := l.list.Accept()
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		wg.Add(1)
		go l.proxy(wg, conn)
	}
	wg.Wait()

	logs.Info("link instance %s shutdown", l.addr)
}

func (l *LinkInstance)Close()  {
	l.Lock()
	l.close = true
	l.list.Close()
	for _, v := range l.channels {
		v.remote.Close()
	}
	l.Unlock()

	l.Wait()
	logs.Info("link instance %s close", l.addr)
}

func (l *LinkInstance)Channels() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.channels)
}

func (l *LinkInstance)Flows() int64 {
	l.RLock()
	defer l.RUnlock()

	var resv, send int64
	for _, v := range l.channels {
		resv += v.resvflow
		send += v.sendflow
	}
	return resv + send
}


func connect(wg *sync.WaitGroup, conn1 net.Conn, conn2 net.Conn, flow *int64)  {
	defer func() {
		wg.Done()
	}()

	var buf [8192]byte
	for  {
		cnt, err1 := conn1.Read(buf[:])
		if cnt > 0 {
			err2 := WriteFull(conn2, buf[:cnt])
			if err2 != nil {
				logs.Error(err2.Error())
				return
			}
			atomic.AddInt64(flow, int64(cnt))
		}
		if err1 != nil {
			logs.Error(err1.Error())
			return
		}
	}
}