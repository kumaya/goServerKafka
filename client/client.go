package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/metadata"

	pb "github.com/kumaya/goServerKafka/proto/manager"
)

const (
	consumerGroup = "sample-client"
)

func HandleConnect(client pb.ManagerClient) error {
	log.Printf("invoked connect to manager")
	ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancelFunc()

	ctx = metadata.AppendToOutgoingContext(ctx, "Realm-ID", "test-consume-once")
	connectStream, err := client.Connect(ctx)
	if err != nil {
		log.Printf("error receiving stream data. err %v", err)
	}
	go func() {
		for {
			log.Printf("sending keepalive message")
			_ = connectStream.Send(&pb.ClientRequest{
				Command: &pb.ClientRequest_Keepalive{Keepalive: 123},
			})
			time.Sleep(20 * time.Second)
		}
	}()

	for {
		chunk, err := connectStream.Recv()
		if err != nil {
			log.Printf("failed while reading: %v", err)
			return err
		}
		if chunk.GetPayload() == nil {
			//log.Print("empty heartbeat message.")
		} else {
			log.Printf("received response: %v", string(chunk.GetPayload()))
			_ = connectStream.Send(&pb.ClientRequest{
				Command: &pb.ClientRequest_Ack{Ack: &pb.ClientRequest_MarkIdentifier{Id: string(chunk.GetPayload())}},
			})
			//log.Printf("received response: %v", chunk.GetTaskID())
		}
	}

}

func HandleOutcome() error {
	panic("not implemented")
}
