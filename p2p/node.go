package p2p

import (
	"math/rand"
	"time"

	"github.com/lokidb/server/pkg/dtypes/dstate"
)

const heartbeatRate = time.Second * 10
const requestsPerBeat = 2
const randomRetrys = 10

type Node interface {
	Start() error
	Shutdown()

	Get(key string) []byte
	Set(key string, value []byte)
	Del(key string)
}

type Address struct {
	Host string
	Port int
}

type node struct {
	server Server
	peers  []Address
	state  dstate.State
	done   chan bool
	run    bool
}

func New(selfAddress Address, bootstrapAddress []Address) Node {
	n := new(node)
	n.state = dstate.New()
	n.peers = bootstrapAddress
	n.server = newNodeServer(selfAddress, &n.state)

	n.run = true
	n.done = make(chan bool, 1)

	return n
}

func (n *node) Start() error {
	if err := n.server.Start(); err != nil {
		return err
	}

	go n.syncLoop()

	return nil
}

func (n *node) Shutdown() {
	n.server.Shutdown()
	n.run = false
	<-n.done
}

func (n *node) Get(key string) []byte {
	return n.state.Get(key)
}

func (n *node) Set(key string, value []byte) {
	n.state.Set(key, value)
}

func (n *node) Del(key string) {
	n.state.Del(key)
}

func (n *node) syncLoop() {
	defer close(n.done)

	for n.run {
		time.Sleep(heartbeatRate)

		requested := make(map[int]struct{})
		for r := 0; r < requestsPerBeat; r++ {
			var randomPeerIndex int = -1

			for i := 0; i < randomRetrys; i++ {
				random := rand.Float64()
				randomPeerIndex = int(random * float64(len(n.peers)))

				if _, ok := requested[randomPeerIndex]; !ok {
					break
				}

				randomPeerIndex = -1
			}

			if randomPeerIndex == -1 {
				continue
			}

			peer := n.peers[randomPeerIndex]
			requested[randomPeerIndex] = struct{}{}

			peerState, err := getStateFromPeer(peer)
			if err != nil {
				continue
			}

			n.state = n.state.Merge(peerState)
		}
	}

	n.done <- true
}
