package main

import (
	"flag"
	"log"
)

var (
	config string
	help   bool
	debug  bool
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.StringVar(&config, "config", "config.yaml", "configure file.")
}

func main() {

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	err := LoadConfig(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	TcpProxyStart()
}
