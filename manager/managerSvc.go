package main

import (
	"context"
	pb "github.com/kumaya/goServerKafka/proto/manager"
	"io"
	"log"
	"sync"
)

var (
	brokerList = []string{"localhost:9092"}
)

const (
	topic    = "test-consume-once"
	clientID = "manager-svc"
)

type managerServer struct {
	pb.ManagerServer
}

func NewManagerServer() pb.ManagerServer {
	return &managerServer{}
}

func (m *managerServer) Connect(stream pb.Manager_ConnectServer) error {
	// TODO:
	//  terminate subscription on client stream close
	//  terminate subscription on server stream close
	//  send stream message to client on message in kafka
	//  ack message in kafka once ack received on client stream.

	parentCtx, cancel := context.WithCancel(stream.Context())
	defer func() {
		log.Printf("context cancel called")
		cancel()
	}()

	//  find the consumerGroup
	consumerGrp := "sample-client"
	log.Printf("connect invoked by %s", consumerGrp)

	//  create kafka subscription for the consumer group
	incomingMessageChan := make(chan interface{})
	clientAckChan := make(chan interface{})
	inFlightMessages := make(map[interface{}]interface{})
	defer func() {
		close(incomingMessageChan)
		close(clientAckChan)
		inFlightMessages = nil
	}()
	kafkaCfg := KakfaConfig{
		ClientID:   clientID,
		BrokerList: brokerList,
		Topic:      []string{topic},
	}
	kafkaConsumer := Consumer{
		Ready:            make(chan bool),
		Once:             sync.Once{},
		Message:          incomingMessageChan,
		InFlightMessages: inFlightMessages,
		Ack:              clientAckChan,
		Mut:              sync.Mutex{},
	}
	err := kafkaCfg.Subscribe(parentCtx, consumerGrp, &kafkaConsumer)
	if err != nil {
		log.Printf("error in kafka, err %v", err)
		return err
	}

	clStreamErr := make(chan error, 1)
	clientReq := make(chan *pb.ClientRequest, 1)
	go func() {
		for {
			clientMsg, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" {
					return
				}
				log.Printf("(managerSvc) (client Recv 1) client stream error: %v", err)
				parentCtx.Done()
				clStreamErr <- err
				return
			}
			select {
			case <-stream.Context().Done():
				parentCtx.Done()
				log.Printf("(managerSvc) (client Recv 2) client stream error: %v", err)
				return
			case clientReq <- clientMsg:
				m := <-clientReq
				switch cmd := m.GetCommand().(type) {
				case *pb.ClientRequest_Keepalive:
					//log.Printf("received keepalive: %d", cmd.Keepalive)
					continue
				case *pb.ClientRequest_Ack:
					clientAckChan <- cmd.Ack.GetId()
				}
			}
		}
	}()

	for {
		log.Printf("inside managerSvc server send for loop")
		select {
		case <-parentCtx.Done():
			log.Printf("(managerSvc) (server send 1) server send error: %v", err)
			return nil
		case err = <-clStreamErr:
			if err == io.EOF {
				log.Printf("(managerSvc) (server send 2) server send error: %v", err)
				return nil
			}
			log.Printf("(managerSvc) (server send 3) server send error: %v", err)
			return err
		case message := <-incomingMessageChan:
			log.Printf("message received from kafka: %s", message)
			err = stream.Send(&pb.ServerResponse{
				Payload: message.([]byte),
			})
			if err != nil {
				log.Printf("error sending payload to client, err %v", err)
				return err
			}
		}

	}
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
