package p2p

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/lokidb/server/pkg/dtypes/dstate"
)

const heartbeatRate = time.Second * 1
const requestsPerBeat = 1

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
	server         Server
	address        Address
	bootstrapPeers []Address
	state          dstate.State
	cancelFunc     context.CancelFunc
}

func New(selfAddress Address, bootstrapAddress []Address) Node {
	n := new(node)
	n.state = dstate.New()
	n.address = selfAddress
	n.bootstrapPeers = bootstrapAddress
	n.server = newNodeServer(selfAddress, n)

	return n
}

func (n *node) Start() error {
	go n.server.Start()

	success := 0
	for _, peer := range n.bootstrapPeers {
		n.addPeerToState(peer, true)

		if peer.Host == n.address.Host && peer.Port == n.address.Port {
			continue
		}

		err := notifyNewPeer(peer, n.address.Host, n.address.Port)
		if err == nil {
			success += 1
		}
	}

	// if success < 1 {
	// 	return fmt.Errorf("did not notify peers about new peer")
	// }

	ctx, cancel := context.WithCancel(context.Background())
	n.cancelFunc = cancel

	go n.syncLoop(ctx)

	return nil
}

func (n *node) Shutdown() {
	n.server.Shutdown()
	n.cancelFunc()
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

func (n *node) syncLoop(ctx context.Context) {
	ticker := time.NewTicker(heartbeatRate)
	for {
		select {
		case <-ticker.C:
			peers := n.peersFromState()
			for r := 0; r < requestsPerBeat; r++ {
				random := rand.Float64()
				randomPeerIndex := int(random * float64(len(peers)))
				peer := peers[randomPeerIndex]

				if peer.Host == n.address.Host && peer.Port == n.address.Port {
					continue
				}

				peerState, err := getStateFromPeer(peer)
				if err != nil {
					continue
				}

				n.state = n.state.Merge(peerState)
			}

		case <-ctx.Done():
			return
		}
	}
}

func (n *node) addPeerToState(peer Address, defualt bool) {
	currentPeers := n.state.Get("peers")
	updatedPeers := append(currentPeers, []byte(fmt.Sprintf("%s:%d,", peer.Host, peer.Port))...)

	if defualt {
		n.state.SetDefault("peers", updatedPeers)
	} else {
		n.state.Set("peers", updatedPeers)
	}
}

func (n *node) peersFromState() []Address {
	peersBytes := n.state.Get("peers")

	peers := make([]Address, 0, 20)
	host := ""
	port := 0
	buff := ""
	for _, b := range peersBytes {
		if b == byte(':') {
			host = buff
			buff = ""
		} else if b == byte(',') {
			port, _ = strconv.Atoi(buff)
			buff = ""
			peers = append(peers, Address{Host: host, Port: port})
		} else {
			buff += string(b)
		}
	}

	return peers
}
