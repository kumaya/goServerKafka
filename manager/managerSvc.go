package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"google.golang.org/grpc/metadata"

	pb "github.com/kumaya/goServerKafka/proto/manager"
)

var (
	brokerList = []string{"localhost:9092"}
)

const (
	clientID = "manager-svc"
)

type managerServer struct {
	pb.ManagerServer
	saramaCfg *sarama.Config
}

func NewManagerServer() pb.ManagerServer {
	// use common sarama config to avoid memory leak so the metrics are registered only once.
	// https://github.com/Shopify/sarama/pull/991#issuecomment-349117911
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.DefaultVersion
	saramaCfg.ClientID = clientID
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}

	return &managerServer{saramaCfg: saramaCfg}
}

func getGroupNTopic(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("%s", "unable to retrieve context metadata.")
	}
	log.Printf("%+v", md)
	entries, ok := md["realm-id"]
	if !ok || len(entries) < 1 {
		return "", fmt.Errorf("unable to find request header 'realm-id' in context metadata")
	}
	last := len(entries) - 1
	return entries[last], nil
}

func (m *managerServer) Connect(stream pb.Manager_ConnectServer) error {
	// Get the consumer group and topic from incoming context
	topic, err := getGroupNTopic(stream.Context())
	if err != nil {
		log.Printf("unable to retrieve metadata from context. err: %v", err)
		return err
	}

	parentCtx, cancel := context.WithCancel(stream.Context())
	defer func() {
		log.Printf("context cancel called")
		cancel()
	}()

	//  find the consumerGroup
	consumerGrp := fmt.Sprintf("%s-%s", "sample-client", topic)
	log.Printf("connect invoked by %s", consumerGrp)

	//  create kafka subscription for the consumer group
	incomingMessageChan := make(chan string)
	clientAckChan := make(chan string)
	var inFlightMessages sync.Map
	defer func() {
		close(incomingMessageChan)
		close(clientAckChan)
		//inFlightMessages = nil
	}()
	kafkaCfg := KakfaConfig{
		ClientID:   clientID,
		BrokerList: brokerList,
		Topic:      []string{topic},
		SaramaCfg:  m.saramaCfg,
	}
	kafkaConsumer := Consumer{
		Ready:            make(chan bool),
		Once:             sync.Once{},
		Message:          incomingMessageChan,
		InFlightMessages: &inFlightMessages,
		Ack:              clientAckChan,
		Mut:              sync.Mutex{},
	}
	err = kafkaCfg.Subscribe(parentCtx, consumerGrp, &kafkaConsumer)
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
				Payload: []byte(message),
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
