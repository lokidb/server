package p2p

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lokidb/server/p2p/state"
)

// Client timeout in seconds
const timeout = 20
const internalKeyPrefix = "$"

type Node interface {
	Run()
	Shutdown()
	UpdateState(key string, value string) error
	OnKeyChange(string, func(string))
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
	handlers      map[string]func(string)
	handledItems  map[string]state.Item
	server        *nodeServer
	stop          bool
	stopChan      chan struct{}
}

func NewNode(peersAddress []Address, heartbeatRate time.Duration, port int) Node {
	n := new(p2pNode)
	n.handledItems = make(map[string]state.Item, 1000)
	n.heartbeatRate = heartbeatRate
	n.stop = false
	n.stopChan = make(chan struct{})
	n.handlers = make(map[string]func(string))
	n.server = newServer(n, "0.0.0.0", port)
	n.state = state.New()

	// Create clients from peers address
	n.clients = make(map[Address]Client, len(peersAddress))
	for _, address := range peersAddress {
		client := newClient(fmt.Sprintf("%s:%d", address.Host, address.Port), time.Duration(time.Second*timeout))
		n.clients[address] = *client
		log.Printf("peer added %d\n", address.Port)
	}

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

		for _, activeItem := range n.state.ActiveItems() {
			currentItem, ok := n.handledItems[activeItem.UID()]
			if ok && currentItem.Value == activeItem.Value {
				continue
			}

			n.handledItems[activeItem.UID()] = activeItem

			handler, ok := n.handlers[activeItem.Key]
			if ok {
				handler(activeItem.Value)
			}
		}
	}
}

func (n *p2pNode) Shutdown() {
	n.stop = true
	n.server.Stop()
	<-n.stopChan
}

func (n *p2pNode) OnKeyChange(key string, handler func(string)) {
	n.handlers[key] = handler
}

func (n *p2pNode) UpdateState(key string, value string) error {
	if strings.HasPrefix(key, internalKeyPrefix) {
		return fmt.Errorf("can't have key with prefix '%s' it reserved for internal node use", internalKeyPrefix)
	}

	n.state.Update(key, value, n.heartbeatRate*3)

	return nil
}

func (n *p2pNode) getState() state.State {
	return n.state
}
