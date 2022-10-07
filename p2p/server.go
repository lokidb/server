package p2p

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lokidb/server/pkg/dtypes/dstate"
	comm "github.com/lokidb/server/pkg/p2p_grpc"
	grpc "google.golang.org/grpc"
)

type Server interface {
	Start() error
	Shutdown()
}

// server is used to implement lokidbService.
type server struct {
	comm.UnimplementedConfigurationNodeServiceServer
	state      *dstate.State
	grpcServer grpc.Server
	host       string
	port       int
}

func newNodeServer(address Address, state *dstate.State) *server {
	s := new(server)
	s.state = state
	grpcServer := grpc.NewServer()
	s.grpcServer = *grpcServer
	s.host = address.Host
	s.port = address.Port

	comm.RegisterConfigurationNodeServiceServer(grpcServer, s)

	return s
}

func (s *server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *server) Shutdown() {
	s.grpcServer.Stop()
}

func (s *server) GetState(ctx context.Context, in *empty.Empty) (*comm.StateResponse, error) {
	state := make(map[string]*comm.Item)
	for key, value := range *(*s.state).GetItems() {
		state[key] = &comm.Item{LastUpdated: value.LastUpdated.UnixNano(), Value: value.Data}
	}

	return &comm.StateResponse{State: state}, nil
}

func (s *server) NewPeer(ctx context.Context, in *comm.NewPeerRequest) (*comm.NewPeerResponse, error) {
	return &comm.NewPeerResponse{Added: true}, nil
}
