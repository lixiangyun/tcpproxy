package main

import (
	"log"
	"net"
	"sync"
	"time"
)

type TcpBridge struct {
	RemoteAddr string
	LocalAddr  string
}

func NewTcpBridge(localaddr string, remoteaddr string) *TcpBridge {
	return &TcpBridge{LocalAddr: localaddr, RemoteAddr: remoteaddr}
}

func tcpBridgeLink(name string, localconn, remoteconn net.Conn) {

	syncSem := new(sync.WaitGroup)

	log.Println("new "+name+" connect. ", localconn.RemoteAddr(), " <-> ", remoteconn.RemoteAddr())

	syncSem.Add(2)

	go tcpChannel(localconn, remoteconn, syncSem)
	go tcpChannel(remoteconn, localconn, syncSem)

	syncSem.Wait()

	log.Println("close "+name+" connect. ", localconn.RemoteAddr(), " <-> ", remoteconn.RemoteAddr())
}

// 桥接连接池模式
func (t *TcpBridge) Bridge(numconn int) error {

	remoteconn := make(chan net.Conn, numconn)
	localconn := make(chan net.Conn, numconn)

	locallisten, err := net.Listen("tcp", t.LocalAddr)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	remotelisten, err := net.Listen("tcp", t.RemoteAddr)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// 监听端口
	go func() {
		for {
			temp, err := locallisten.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			localconn <- temp
		}
	}()

	// 监听端口
	go func() {
		for {
			temp, err := remotelisten.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			remoteconn <- temp
		}
	}()

	// 等待连接建立
	for {
		conn1 := <-remoteconn
		conn2 := <-localconn

		go tcpBridgeLink("bridge", conn1, conn2)
	}

	return nil
}

// 单次连接模式
func (t *TcpBridge) LinkOnce() error {

	localconn, err := net.Dial("tcp", t.LocalAddr)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	remoteconn, err := net.Dial("tcp", t.RemoteAddr)
	if err != nil {
		log.Println(err.Error())
		localconn.Close()
		return err
	}

	tcpBridgeLink("linkonce", localconn, remoteconn)

	return nil
}

func tcpBridgeLinkPools(conn1, conn2 net.Conn, sem chan struct{}) {
	tcpBridgeLink("linkpools", conn1, conn2)
	sem <- struct{}{}
}

// 连接池模式
func (t *TcpBridge) LinkPools(numconn int) {

	log.Println("link pools num : ", numconn)

	sem := make(chan struct{}, numconn)

	for i := 0; i < numconn; i++ {
		sem <- struct{}{}
	}

	for {

		<-sem

		localconn, err := net.Dial("tcp", t.LocalAddr)
		if err != nil {
			log.Println(err.Error())

			log.Println("retry to link ", t.LocalAddr, " <-> ", t.RemoteAddr)

			time.Sleep(5 * time.Second)

			sem <- struct{}{}
			continue
		}

		remoteconn, err := net.Dial("tcp", t.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			localconn.Close()

			log.Println("retry to link ", t.LocalAddr, " <-> ", t.RemoteAddr)

			time.Sleep(5 * time.Second)

			sem <- struct{}{}
			continue
		}

		go tcpBridgeLinkPools(localconn, remoteconn, sem)
	}
}
