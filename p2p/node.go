package p2p

import (
	"fmt"
	"log"
	"time"

	"github.com/lokidb/server/p2p/state"
)

// Client timeout in seconds
const timeout = 20

type Node interface {
	Run()
	Shutdown()
	UpdateState(key string, value any)
	OnKeyChange(string, func(any))
	getState() state.State
}

type Address struct {
	Host string
	Port int
}

type p2pNode struct {
	clients       map[Address]Client
	heartbeatRate time.Duration
	state         state.State
	handlers      map[string]map[string]func(any)
	server        *nodeServer
	stop          bool
	stopChan      chan struct{}
}

func NewNode(peersAddress []Address, heartbeatRate time.Duration, port int) Node {
	n := new(p2pNode)
	n.heartbeatRate = heartbeatRate
	n.stop = false
	n.stopChan = make(chan struct{})
	n.handlers = make(map[string]map[string]func(any))

	// Create clients from peers address
	n.clients = make(map[Address]Client, len(peersAddress))
	for _, address := range peersAddress {
		client := newClient(fmt.Sprintf("%s:%d", address.Host, address.Port), time.Duration(time.Second*timeout))
		n.clients[address] = *client
		log.Printf("peer added %d\n", address.Port)
	}

	n.server = newServer(n, "0.0.0.0", port)

	n.state = state.New()
	n.state.Update("internal", "peers", peersAddress, n.heartbeatRate*3)

	return n
}

// Run node main loop.
func (n *p2pNode) Run() {
	defer close(n.stopChan)

	// Run server
	go n.server.Run()

	// Run until shutdown
	for !n.stop {
		time.Sleep(n.heartbeatRate)

		// get peers state and merge to self state
		for _, client := range n.clients {
			clientState, err := client.GetState()
			if err != nil {
				continue
			}

			log.Println("Merging state from peer")
			n.state.Merge(clientState)
		}
	}
}

func (n *p2pNode) Shutdown() {
	n.stop = true
	n.server.Stop()
	<-n.stopChan
}

func (n *p2pNode) addHandler(name string, key string, handler func(any)) {
	n.handlers[name][key] = handler
}

func (n *p2pNode) OnKeyChange(key string, handler func(any)) {
	n.addHandler("exteranl", key, handler)
}

func (n *p2pNode) UpdateState(key string, value any) {
	n.state.Update("external", key, value, n.heartbeatRate*3)
}

func (n *p2pNode) getState() state.State {
	return n.state
}
