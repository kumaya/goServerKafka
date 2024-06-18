package main

import (
	"context"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

type KakfaConfig struct {
	ClientID   string
	BrokerList []string
	Topic      []string
	SaramaCfg  *sarama.Config
}

type Consumer struct {
	Ready            chan bool
	Once             sync.Once
	Message          chan<- string
	InFlightMessages *sync.Map
	Ack              <-chan string
	Mut              sync.Mutex
}

func (k *KakfaConfig) Subscribe(ctx context.Context, groupId string, consumer *Consumer) error {
	consumerGroup, err := sarama.NewConsumerGroup(k.BrokerList, groupId, k.SaramaCfg)
	if err != nil {
		return err
	}

	go func() {
		for {
			if ctx.Err() != nil {
				log.Printf("(kafkaConsumer) client connection terminated for context error. err %v", ctx.Err())
				return
			}
			log.Printf("starting loop to consume")
			if err = consumerGroup.Consume(ctx, k.Topic, consumer); err != nil {
				log.Printf("error while starting consume loop: %v", err)
				return
			}
		}
	}()
	<-consumer.Ready
	log.Printf("sarama consumer up and running. Consuming messages...")
	//wg.Wait()
	return nil
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	consumer.Once.Do(func() {
		log.Printf("closing ready channel")
		close(consumer.Ready)
	})
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			consumer.Message <- string(message.Value)
			consumer.InFlightMessages.Store(string(message.Value), message)
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			//session.MarkMessage(message, "")
		case <-session.Context().Done():
			log.Printf("(kafkaConsumer/ConsumeClaim) closing")
			return nil
		case messageStr := <-consumer.Ack:
			log.Printf("message acked; %s %v", messageStr, consumer.InFlightMessages)
			kMsg, ok := consumer.InFlightMessages.Load(messageStr)
			if !ok {
				log.Printf("not found %s", messageStr)
				continue
			}
			consumer.InFlightMessages.Delete(messageStr)
			session.MarkMessage(kMsg.(*sarama.ConsumerMessage), "")
		}
	}
}
