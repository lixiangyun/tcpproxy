package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

var gtotalUpSize uint64
var gtotalDownSize uint64

func init() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-ticker.C
			display()
		}
	}()
}

func Add(up int, down int) {
	atomic.AddUint64(&gtotalUpSize, uint64(up))
	atomic.AddUint64(&gtotalDownSize, uint64(down))
}

func display() {
	log.Printf("↑%s ↓%s\n",
		calcUnit(gtotalUpSize), calcUnit(gtotalDownSize))
}

func calcUnit(cnt uint64) string {
	if cnt < 1024 {
		return fmt.Sprintf("%d", cnt)
	} else if cnt < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float32(cnt)/1024)
	} else if cnt < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float32(cnt)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float32(cnt)/(1024*1024*1024))
	}
}
