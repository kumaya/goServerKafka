package main

import (
	"context"
	"log"
	"time"

	pb "github.com/kumaya/goServerKafka/proto/manager"
)

type managerServer struct {
	pb.ManagerServer
}

func NewManagerServer() pb.ManagerServer {
	return &managerServer{}
}

func (m *managerServer) Connect(request *pb.ConnectRequest, server pb.Manager_ConnectServer) error {
	cChan := make(chan interface{})
	log.Printf("connect invoked by %s", request.GetConsumerGroup())
	if err := ConnectAndConsume(server.Context(), request.GetConsumerGroup(), cChan); err != nil {
		log.Printf("error in kafka, err %v", err)
		return err
	}
	for {
		select {
		case message := <-cChan:
			log.Printf("message received from kafka: %s", message)
			err := server.Send(&pb.ConnectResponse{
				Payload: message.([]byte),
			})
			if err != nil {
				log.Printf("error sending payload to client, err %v", err)
				return err
			}
		case <-server.Context().Done():
			log.Printf("client closed the connection")
			return nil
		default:
			time.Sleep(1 * time.Second)
			err := server.Send(&pb.ConnectResponse{})
			if err != nil {
				log.Printf("error sending heartbeat to client, err %v", err)
				return err
			}
		}

	}
	return nil
}

func (m *managerServer) Outcome(ctx context.Context, request *pb.OutcomeRequest) (*pb.OutcomeResponse, error) {
	log.Printf("invoked outcome for task: %d", request.GetTaskID())
	log.Printf("incoming payload: %+v", request.GetResponsePayload())
	return &pb.OutcomeResponse{}, nil
}

func (m *managerServer) mustEmbedUnimplementedManagerServer() {
	//TODO implement me
	panic("implement me")
}
