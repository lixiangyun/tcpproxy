package main

import (
	"crypto/tls"
	"fmt"

	//"io"
	"log"
	"net"
	"sync"
	"time"
)

type TcpProxy struct {
	ListenTls  *tls.Config
	ListenAddr string
	RemoteTls  *tls.Config
	RemoteAddr []string
}

func NewTcpProxy(local string, localtls *tls.Config, remote []string, remotetls *tls.Config) *TcpProxy {
	return &TcpProxy{ListenTls: localtls, ListenAddr: local, RemoteTls: remotetls, RemoteAddr: remote}
}

func writeFull(conn net.Conn, buf []byte) error {
	totallen := len(buf)
	sendcnt := 0

	for {
		cnt, err := conn.Write(buf[sendcnt:])
		if err != nil {
			return err
		}
		if cnt+sendcnt >= totallen {
			return nil
		}
		sendcnt += cnt
	}
}

// tcp通道互通
func tcpChannel(prefix string, localconn net.Conn, remoteconn net.Conn, wait *sync.WaitGroup) {
	defer wait.Done()
	defer localconn.Close()
	defer remoteconn.Close()
	buf := make([]byte, 65535)
	for {
		cnt, err := localconn.Read(buf[0:])
		if err != nil {
			if cnt != 0 {
				writeFull(remoteconn, buf[0:cnt])
			}
			break
		}
		if debug {
			log.Printf("%s body:[%v]\r\n", prefix, buf[0:cnt])
		}
		err = writeFull(remoteconn, buf[0:cnt])
		if err != nil {
			break
		}
	}
}

// tcp代理处理
func tcpProxyProcess(localconn net.Conn, remoteconn net.Conn) {

	localremote := fmt.Sprintf("%s->%s",
		localconn.RemoteAddr().String(),
		remoteconn.RemoteAddr().String())

	remotelocal := fmt.Sprintf("%s->%s",
		remoteconn.RemoteAddr().String(),
		localconn.RemoteAddr().String())

	log.Println("new connect. ", localremote)

	syncSem := new(sync.WaitGroup)
	syncSem.Add(2)
	go tcpChannel(localremote, localconn, remoteconn, syncSem)
	go tcpChannel(remotelocal, remoteconn, localconn, syncSem)
	syncSem.Wait()

	log.Println("close connect. ", localremote)
}

// 正向tcp代理启动和处理入口
func (t *TcpProxy) Start() error {
	var times int

	listen, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	var remoteaddr string
	for _, v := range t.RemoteAddr {
		remoteaddr += v + " "
	}

	log.Printf("listen : %s -> %s", t.ListenAddr, remoteaddr)

	for {
		var localconn net.Conn
		var remoteconn net.Conn

		localconn, err = listen.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if t.ListenTls != nil {
			localconn = tls.Server(localconn, t.ListenTls)
		}

		for i := 0; i < len(t.RemoteAddr); i++ {
			remoteaddr := t.RemoteAddr[times]
			times = (times + 1) % len(t.RemoteAddr)

			remoteconn, err = net.Dial("tcp", remoteaddr)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			log.Println("proxy connect to ", remoteaddr)
			break
		}

		if remoteconn == nil {
			localconn.Close()
			continue
		}

		if t.RemoteTls != nil {
			remoteconn = tls.Client(remoteconn, t.RemoteTls)
		}

		go tcpProxyProcess(localconn, remoteconn)
	}

	return nil
}

func TcpProxyStart() {

	listeners := listenerGetAll()
	if 0 == len(listeners) {
		log.Fatalln("no listenner.")
	}

	for _, v := range listeners {

		var localtls *tls.Config
		var remotetls *tls.Config

		tls := TlsGet(v.Tlsname)
		if tls != nil {
			localtls = TlsServerConfig(tls)
		}

		cluster := ClusterGet(v.Cluster)
		if cluster == nil {
			log.Fatalf("not found %s cluster.", v.Cluster)
		}

		if len(cluster.Endpoint) == 0 {
			log.Fatalf("not found %s cluster endpoint.", v.Cluster)
		}

		tls = TlsGet(cluster.TlsName)
		if tls != nil {
			remotetls = TlsClientConfig(tls, cluster.Endpoint[0])
		}

		tcoporxy := NewTcpProxy(v.Address, localtls, cluster.Endpoint, remotetls)

		go func() {
			err := tcoporxy.Start()
			if err != nil {
				log.Fatalf("tcp proxy start failed %v.", v)
			}
		}()
	}

	for {
		time.Sleep(time.Second * 100)
	}
}
