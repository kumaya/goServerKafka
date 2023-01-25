# goServerKafka

- creates a grpc manager service and a grpc client service
- The job of the manager:
  - expose a stream endpoint
  - send periodic heartbeat to keep connection alive.
  - create kafka consumer based on incoming request.ConsumerGroup
  - consume from the kafka topic
  - send the message to one of the connected clients only having similar consumerGroup.
- The job of the client:
  - connect to manager service
  - consume connect message
  - run multiple instances sharing same ConsumerGroup
  - only one client should log the message.

### Useful Commands
```
kafka-topics --bootstrap-server broker:9092 --create --topic test-consume-once --partitions 3
kafka-topics --bootstrap-server broker:9092 --describe --topic test-consume-once

kafka-consumer-groups  --bootstrap-server broker:9092 --list
kafka-consumer-groups  --bootstrap-server broker:9092 --describe --group sample-client

produce message
kafka-console-producer --bootstrap-server broker:9092 --topic test-consume-once
```