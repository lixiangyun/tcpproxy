package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"
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
	buf := make([]byte, 4096)
	for {
		cnt, err := localconn.Read(buf[0:])
		if err != nil {
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
			log.Println("proxy connect to ", t.RemoteAddr[times])
			remoteconn, err = net.Dial("tcp", t.RemoteAddr[times])
			if err != nil {
				log.Println(err.Error())
				times++
				if times >= len(t.RemoteAddr) {
					times = 0
				}
				log.Printf("switch address to %s.\n", t.RemoteAddr[times])
			} else {
				break
			}
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
			remotetls = TlsClientConfig(tls)
		}

		tcoporxy := NewTcpProxy(v.Address, localtls, cluster.Endpoint, remotetls)
		err := tcoporxy.Start()
		if err != nil {
			log.Fatalf("tcp proxy start failed %v.", v)
		}
	}

}
