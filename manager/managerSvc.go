package main

import (
	"context"
	pb "github.com/kumaya/goServerKafka/proto/manager"
	"log"
	"time"
)

type managerServer struct {
	pb.ManagerServer
}

func NewManagerServer() pb.ManagerServer {
	return &managerServer{}
}

func (m *managerServer) Connect(request *pb.ConnectRequest, server pb.Manager_ConnectServer) error {
	log.Printf("connect invoked by %s", request.GetConsumerGroup())
	for {
		select {
		default:
			time.Sleep(1 * time.Second)
			err := server.Send(&pb.ConnectResponse{})
			if err != nil {
				log.Printf("error sending to client, err %v", err)
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
