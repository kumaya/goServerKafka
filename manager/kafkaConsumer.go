package main

import (
	"context"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	brokerList = []string{"localhost:9092"}
)

const (
	topic    = "test-consume-once"
	clientID = "manager-svc"
)

func ConnectAndConsume(ctx context.Context, cg string, consumeChan chan<- interface{}) error {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.ClientID = clientID
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}

	consumer := Consumer{ready: make(chan bool), msg: consumeChan}
	consumerGroup, err := sarama.NewConsumerGroup(brokerList, cg, config)
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Printf("error from consumer: %v", err)
				return
			}
			if ctx.Err() != nil {
				log.Printf("client connection terminated. err %v", err)
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	<-consumer.ready
	log.Printf("sarama consumer up and running")
	//wg.Wait()
	log.Printf("consuming message...")
	return nil
}

type Consumer struct {
	ready chan bool
	msg   chan<- interface{}
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			consumer.msg <- message.Value
			//log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
