package client

import (
	"context"
	"log"
	"time"

	pd "github.com/lokidb/server/communication/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	grpcConn   grpc.ClientConn
	grpcClient pd.LokiDBServiceClient
	timeout    time.Duration
}

func New(addr string, timeout time.Duration) *Client {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := new(Client)
	client.grpcClient = pd.NewLokiDBServiceClient(conn)
	client.grpcConn = *conn
	client.timeout = timeout

	return client
}

func (c *Client) Close() {
	c.grpcConn.Close()
}

func (c *Client) Get(key string) (string, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.grpcClient.Get(ctx, &pd.GetRequest{Key: key})

	if err != nil {
		return "", err
	}

	return res.GetValue(), nil
}

func (c *Client) Set(key string, value string) error {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.grpcClient.Set(ctx, &pd.SetRequest{Key: key, Value: value})

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Del(key string) (bool, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.grpcClient.Del(ctx, &pd.DelRequest{Key: key})
	return res.GetDeleted(), err
}

func (c *Client) Keys() ([]string, error) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.grpcClient.Keys(ctx, &emptypb.Empty{})

	return res.GetKeys(), err
}

func (c *Client) Flush() error {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.grpcClient.Flush(ctx, &emptypb.Empty{})

	return err
}
