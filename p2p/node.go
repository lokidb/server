package p2p

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// Client timeout in seconds
const timeout = 20

type Node interface {
	Run()
	Shutdown()
	SendMessage(message)
	OnMessage(msgName string, handler func(payload string))
	GetState() State
}

type Address struct {
	Host string
	Port int
}

type p2pNode struct {
	uuid          string
	peers         []Address
	clients       []*Client
	state         State
	handlers      map[string]func(string)
	heartbeatRate time.Duration
	stop          bool
	stopChan      chan struct{}
	server        *nodeServer
}

func NewNode(peers []Address, heartbeatRate time.Duration) Node {
	n := new(p2pNode)
	n.uuid = uuid.New().String()
	n.peers = peers
	n.clients = make([]*Client, 0, len(peers))
	n.handlers = make(map[string]func(string), 20)
	n.state = newState()
	n.stop = false
	n.stopChan = make(chan struct{})
	n.server = newServer(n, "0.0.0.0", 11497)

	for _, peer := range peers {
		client := newClient(fmt.Sprintf("%s:%d", peer.Host, peer.Port), time.Duration(time.Second*timeout))
		n.clients = append(n.clients, client)
	}

	return n
}

func (n *p2pNode) Run() {
	defer close(n.stopChan)

	for !n.stop {
		time.Sleep(n.heartbeatRate)

		// get peers state and merge to self state
		for _, client := range n.clients {
			clientState, err := client.GetState()
			if err != nil {
				log.Println(err.Error())
				continue
			}

			n.state = n.state.merge(clientState)
		}

		// sign on state
		n.state.sign(n.uuid)

		// notify handlers
		for _, msg := range n.state.messages {
			if msg.handled {
				continue
			}

			if len(msg.signaturs) == len(n.peers) {
				handler, ok := n.handlers[msg.name]

				if ok {
					handler(msg.payload)
					msg.handled = true
				}
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

func (n *p2pNode) GetState() State {
	return n.state
}
