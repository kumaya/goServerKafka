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
