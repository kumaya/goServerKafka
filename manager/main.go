package main

import (
	pb "github.com/kumaya/goServerKafka/proto/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"time"
)

const (
	port = ":8080"
)

func main() {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    30 * time.Second,
				Timeout: 60 * time.Second,
			},
		))
	mgrSvr := NewManagerServer()
	pb.RegisterManagerServer(grpcServer, mgrSvr)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s, err: %v", port, err)
	}
	log.Printf("Manager server started. Listening on port: %s", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	defer grpcServer.GracefulStop()
	return
}

// create kafka consumer from incoming connection consumer group
// consume msg from kafka and give to an active member of a consumer group
