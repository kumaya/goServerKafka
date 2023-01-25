package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

var (
	brokerList = []string{"localhost:9092"}
)

const (
	topic    = "test-consume-once"
	clientID = "manager-svc"
)

func main() {
	pConf := sarama.NewConfig()
	pConf.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	pConf.Producer.RequiredAcks = sarama.WaitForAll
	pConf.Producer.Return.Successes = true
	pConf.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokerList, pConf)
	if err != nil {
		log.Panicf("Error initializing producer: %v", err)
	}
	log.Printf("starting new kafka producer")
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()
	for count := 0; count <= 60; count++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("kuchBhi %d", count)),
		}
		p, o, err := producer.SendMessage(msg)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, p, o)
		time.Sleep(2 * time.Second)
	}
}
