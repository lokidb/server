package p2p

import (
	"fmt"
	"log"
	"time"
)

// Client timeout in seconds
const timeout = 20

type Node interface {
	Run()
	Shutdown()
	SendMessage(message)
	OnMessage(msgName string, handler func(payload string))
	getState() State
}

type Address struct {
	Host string
	Port int
}

type p2pNode struct {
	peers         map[Address]Client
	state         State
	handlers      map[string]func(string)
	heartbeatRate time.Duration
	stop          bool
	stopChan      chan struct{}
	server        *nodeServer
}

func NewNode(peersAddress []Address, heartbeatRate time.Duration, port int) Node {
	n := new(p2pNode)
	n.peers = make(map[Address]Client, len(peersAddress))

	// Create clients from peers address
	for _, address := range peersAddress {
		client := newClient(fmt.Sprintf("%s:%d", address.Host, address.Port), time.Duration(time.Second*timeout))
		n.peers[address] = *client
		log.Printf("peer added %d\n", address.Port)
	}

	n.server = newServer(n, "0.0.0.0", port)
	n.handlers = make(map[string]func(string), 20)
	n.state = newState()

	n.heartbeatRate = heartbeatRate
	n.stop = false
	n.stopChan = make(chan struct{})

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
		for _, client := range n.peers {
			clientState, err := client.GetState()
			if err != nil {
				continue
			}

			log.Println("Merging state from peer")
			n.state = n.state.merge(clientState)
		}

		// notify handlers
		for _, msg := range n.state.messages {
			if !msg.isActive() {
				handler, ok := n.handlers[msg.name]

				if ok {
					handler(msg.payload)
				}

				delete(n.state.messages, msg.id)
			}
		}
	}
}

func (n *p2pNode) Shutdown() {
	n.stop = true
	<-n.stopChan
}

func (n *p2pNode) SendMessage(msg message) {
	n.state.AddMessage(msg)
}

func (n *p2pNode) OnMessage(msgName string, handler func(payload string)) {
	n.handlers[msgName] = handler
}

func (n *p2pNode) getState() State {
	return n.state
}
