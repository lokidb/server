package main

import (
	"fmt"
	"time"

	"github.com/lokidb/server/p2p"
)

func main() {
	BootstrapPeers := []p2p.Address{
		{Host: "127.0.0.1", Port: 5212},
		{Host: "127.0.0.1", Port: 5887},
	}
	node := p2p.NewNode(BootstrapPeers, time.Second*5)

	node.Run()
	defer node.Shutdown()

	node.OnMessage("hello-world", func(payload string) {
		fmt.Println(payload)
	})
}
