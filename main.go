package main

import (
	"log"
	"net"

	"github.com/ez-deploy/cluster/service"
	pb "github.com/ez-deploy/protobuf/project"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:80")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc, err := service.NewInClusterService()
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterProjectOpsServer(s, svc)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
