package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/lokidb/server/p2p"
)

var id = flag.Int("id", 0, "node id")

func main() {
	flag.Parse()

	var host string
	var port int
	var bootstrapHost string = "127.0.0.1"
	var bootstrapPort int = 1111

	switch *id {
	case 0:
		host = "127.0.0.1"
		port = 1111
	case 1:
		host = "127.0.0.2"
		port = 2222
	case 2:
		host = "127.0.0.3"
		port = 3333

	}

	node := p2p.New(p2p.Address{Host: host, Port: port}, []p2p.Address{{Host: bootstrapHost, Port: bootstrapPort}})

	node.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ticker := time.NewTicker(time.Second * 2)
	stop := false
	for !stop {
		select {
		case <-c:
			node.Shutdown()
			stop = true
		case <-ticker.C:
			count := node.Get("count")
			if count == nil {
				count = []byte{0}
			}
			log.Printf("count %d\n", count[0])
			count[0] += 1
			node.Set("count", count)
		}
	}
	fmt.Println("goodbay")
}
