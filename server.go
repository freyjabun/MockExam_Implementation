package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "example.com/increment"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedIncrementServiceServer
	inc incrementor
}

type incrementor struct {
	lock         sync.Mutex
	currentValue int32
}

func main() {
	ownPort := ":5000"

	log.Printf("Listening on port %v", ownPort)
	lis, err := net.Listen("tcp", ownPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %v, %v", ownPort, err)
	}

	grpcServer := grpc.NewServer()

	s := &server{inc: incrementor{currentValue: 0}}

	pb.RegisterIncrementServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func (s *server) Increment(ctx context.Context, req *pb.IncrementRequest) (*pb.IncrementReply, error) {
	s.inc.lock.Lock()
	defer s.inc.lock.Unlock()

	rep := &pb.IncrementReply{ValueBefore: s.inc.currentValue}

	if req.Value > s.inc.currentValue {
		rep.Success = true
		s.inc.currentValue = req.Value
	} else {
		rep.Success = false
	}
	return rep, nil
}
