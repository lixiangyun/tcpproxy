package main

import (
	"io"
	"log"
	"net"
	"sync"
)

type TcpProxy struct {
	ListenAddr string
	RemoteAddr string
}

func NewTcpProxy(local string, remote string) *TcpProxy {
	return &TcpProxy{ListenAddr: local, RemoteAddr: remote}
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
func tcpChannel(localconn net.Conn, remoteconn net.Conn, wait *sync.WaitGroup) {

	defer wait.Done()
	defer localconn.Close()
	defer remoteconn.Close()

	buf := make([]byte, 4096)

	for {
		cnt, err := localconn.Read(buf[0:])
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
			}
			break
		}

		err = writeFull(remoteconn, buf[0:cnt])
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
			}
			break
		}
	}
}

// tcp代理处理
func tcpProxyProcess(localconn net.Conn, remoteconn net.Conn) {

	syncSem := new(sync.WaitGroup)

	log.Println("new connect. ", localconn.RemoteAddr(), " -> ", remoteconn.RemoteAddr())

	syncSem.Add(2)

	go tcpChannel(localconn, remoteconn, syncSem)
	go tcpChannel(remoteconn, localconn, syncSem)

	syncSem.Wait()

	log.Println("close connect. ", localconn.RemoteAddr(), " -> ", remoteconn.RemoteAddr())
}

// 正向tcp代理启动和处理入口
func (t *TcpProxy) Start() error {

	listen, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	for {
		localconn, err := listen.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		remoteconn, err := net.Dial("tcp", t.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			localconn.Close()
			continue
		}

		go tcpProxyProcess(localconn, remoteconn)
	}

	return nil
}
