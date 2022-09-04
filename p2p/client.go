package p2p

import (
	"context"
	"log"
	"time"

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
		log.Fatalf("did not connect: %v", err)
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

func (c *Client) GetState() (State, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.grpcClient.GetState(ctx, &emptypb.Empty{})

	if err != nil {
		return State{}, err
	}

	s := newState()

	for _, stateMsg := range res.GetMessages() {
		msg := message{
			id:      stateMsg.Id,
			name:    stateMsg.Name,
			payload: stateMsg.Payload,
			created: time.Unix(stateMsg.Created, 0),
			maxLife: time.Duration(stateMsg.MaxLife * int64(time.Second)),
		}
		s.AddMessage(msg)
	}

	return s, nil
}
