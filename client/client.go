package main

import (
	"context"
	pb "github.com/kumaya/goServerKafka/proto/manager"
	"log"
)

const (
	consumerGroup = "sample-client"
)

func HandleConnect(client pb.ManagerClient) error {
	log.Printf("invoked connect to manager")
	ctx := context.Background()
	connectStream, err := client.Connect(ctx, &pb.ConnectRequest{ConsumerGroup: consumerGroup})
	if err != nil {
		log.Printf("error receiving stream data. err %v", err)
	}
	for {
		chunk, err := connectStream.Recv()
		if err != nil {
			log.Printf("failed while reading: %v", err)
			return err
		}
		if chunk.GetPayload() == nil || chunk.GetTaskID() == 0 {
			log.Print("empty heartbeat message.")
		} else {
			log.Printf("received response: %v", chunk.GetPayload())
		}
	}

}

func HandleOutcome() error {
	panic("not implemented")
}
