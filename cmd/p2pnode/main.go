package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lokidb/server/p2p"
)

var bootstarp = flag.Int("bootstrap", 5544, "Bootstrap port")
var port = flag.Int("port", 4455, "Node port")

func main() {
	flag.Parse()

	BootstrapPeers := []p2p.Address{
		{Host: "127.0.0.1", Port: *bootstarp},
	}
	node := p2p.NewNode(BootstrapPeers, time.Second*2, *port)

	node.SendMessage(p2p.NewMessage(fmt.Sprintf("hello-from-%d", *port), "hello-world", fmt.Sprintf("hello %d", *bootstarp), time.Second*10))

	node.OnMessage("hello-world", func(payload string) {
		fmt.Println(payload)
	})

	node.Run()
	defer node.Shutdown()
}
