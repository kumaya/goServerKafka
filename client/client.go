package main

import (
	"context"
	"log"
	"time"

	pb "github.com/kumaya/goServerKafka/proto/manager"
)

const (
	consumerGroup = "sample-client"
)

func HandleConnect(client pb.ManagerClient) error {
	log.Printf("invoked connect to manager")
	ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))
	defer cancelFunc()
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
		if chunk.GetPayload() == nil {
			//log.Print("empty heartbeat message.")
		} else {
			log.Printf("received response: %v", string(chunk.GetPayload()))
			//log.Printf("received response: %v", chunk.GetTaskID())
		}
	}

}

func HandleOutcome() error {
	panic("not implemented")
}
