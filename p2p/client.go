package p2p

import (
	"context"
	"log"
	"time"

	"github.com/lokidb/server/p2p/state"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	grpcConn   grpc.ClientConn
	grpcClient P2PServiceClient
	timeout    time.Duration
}

func newClient(addr string, timeout time.Duration) *Client {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v\n", err)
	}

	client := new(Client)
	client.grpcClient = NewP2PServiceClient(conn)
	client.grpcConn = *conn
	client.timeout = timeout

	return client
}

func (c *Client) Close() {
	c.grpcConn.Close()
}

func (c *Client) GetState() (state.State, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.grpcClient.GetState(ctx, &emptypb.Empty{})

	if err != nil {
		return state.State{}, err
	}

	s := state.New()

	for _, item := range res.GetStateItems() {
		s.Update(item.GetKey(), item.GetValue(), time.Duration(item.GetInActiveDuration()))
	}

	return s, nil
}
