package grpc

import (
	context "context"
	"fmt"
	"log"
	"net"

	"github.com/lokidb/engine"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	UnimplementedLokiDBServiceServer
	engine     engine.KeyValueStore
	grpcServer grpc.Server
	host       string
	port       int
}

func NewServer(engine *engine.KeyValueStore, host string, port int) *server {
	grpcServer := grpc.NewServer()
	s := new(server)
	s.engine = *engine
	s.grpcServer = *grpcServer
	s.host = host
	s.port = port

	RegisterLokiDBServiceServer(grpcServer, s)

	return s
}

func (s *server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Get(ctx context.Context, in *GetRequest) (*GetResponse, error) {
	value := s.engine.Get(in.GetKey(), nil)

	return &GetResponse{Value: string(value)}, nil
}

func (s *server) Set(ctx context.Context, in *SetRequest) (*emptypb.Empty, error) {
	err := s.engine.Set(in.GetKey(), []byte(in.GetValue()))
	return &emptypb.Empty{}, err
}

func (s *server) Del(ctx context.Context, in *DelRequest) (*DelResponse, error) {
	deleted := s.engine.Del(in.GetKey())

	return &DelResponse{Deleted: deleted}, nil
}

func (s *server) Keys(ctx context.Context, in *emptypb.Empty) (*KeysResponse, error) {
	keys := s.engine.Keys()

	return &KeysResponse{Keys: keys}, nil
}

func (s *server) Flush(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	s.engine.Flush()

	return &emptypb.Empty{}, nil
}
