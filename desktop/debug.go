package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
)

var (
	debug bool
	port  int
)

func init()  {
	flag.BoolVar(&debug, "debug", false, "running pprof for golang debug.")
	flag.IntVar(&port, "port", 18000, "debug pprof http server listen.")
}

func DebugFlag() bool {
	return debug
}

func DebugInit() error {
	flag.Parse()
	if !debug {
		return nil
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	pprofHandler := http.NewServeMux()
	pprofHandler.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	server := &http.Server{Handler: pprofHandler}

	go func() {
		err := server.Serve(listen)
		if err != nil {
			logs.Error(err.Error())
		}
	}()

	return nil
}