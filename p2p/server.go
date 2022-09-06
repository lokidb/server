package p2p

import (
	context "context"
	"fmt"
	"log"
	"net"

	grpclib "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type nodeServer struct {
	UnimplementedP2PServiceServer
	grpcServer grpclib.Server
	host       string
	port       int
	p2pNode    Node
}

func newServer(node Node, host string, port int) *nodeServer {
	grpcServer := grpclib.NewServer()
	s := new(nodeServer)
	s.grpcServer = *grpcServer
	s.host = host
	s.port = port
	s.p2pNode = node

	RegisterP2PServiceServer(grpcServer, s)

	return s
}

func (s *nodeServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("P2P node listening at %v", lis.Addr())
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *nodeServer) GetState(ctx context.Context, in *emptypb.Empty) (*GetStateResponse, error) {
	state := s.p2pNode.getState()
	messages := state.messages
	stateMessages := make(map[string]*StateMessage, len(messages))

	for _, msg := range messages {
		stateMessages[msg.id] = &StateMessage{Id: msg.id, Name: msg.name, Payload: msg.payload, Created: msg.created.Unix(), MaxLife: int64(msg.maxLife.Seconds())}
	}

	return &GetStateResponse{Messages: stateMessages}, nil
}
