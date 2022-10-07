package p2p

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lokidb/server/pkg/dtypes/dstate"
	pd "github.com/lokidb/server/pkg/p2p_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	grpcConn   grpc.ClientConn
	grpcClient pd.ConfigurationNodeServiceClient
	timeout    time.Duration
}

func newClient(addr string, timeout time.Duration) *Client {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := new(Client)
	client.grpcClient = pd.NewConfigurationNodeServiceClient(conn)
	client.grpcConn = *conn
	client.timeout = timeout

	return client
}

func (c *Client) Close() {
	c.grpcConn.Close()
}

func (c *Client) getState() (dstate.State, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.grpcClient.GetState(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	state := dstate.New()
	for key, value := range res.GetState() {
		(*state.GetItems())[key] = dstate.Value{LastUpdated: time.UnixMicro(value.GetLastUpdated()), Data: value.Value}
	}

	return state, nil
}

func (c *Client) NewPeer(host string, port int) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.grpcClient.NewPeer(ctx, &pd.NewPeerRequest{Host: host, Port: int64(port)})

	if err != nil {
		return err
	}

	if !resp.GetAdded() {
		return fmt.Errorf("peer not added")
	}

	return nil
}

func getStateFromPeer(peer Address) (dstate.State, error) {
	client := newClient(fmt.Sprintf("%s:%d", peer.Host, peer.Port), time.Second*2)
	return client.getState()
}
