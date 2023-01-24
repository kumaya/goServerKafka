package main

import (
	"context"
	pb "github.com/kumaya/goServerKafka/proto/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	mgrSvrAddress = "localhost:8080"
	ctxTimeout    = 60 * time.Second
)

func main() {
	// dial to server
	ctx, cancelFunc := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancelFunc()
	conn, err := grpc.DialContext(
		ctx,
		mgrSvrAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to listen to manager server; err %v", err)
	}
	defer conn.Close()

	mgrClient := pb.NewManagerClient(conn)
	err = HandleConnect(mgrClient)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return
}
