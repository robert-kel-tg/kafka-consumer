# kafka web UI
http://localhost:9021

# kafka-consumer

## Create topic
docker-compose exec kafka  \
kafka-topics --create --topic foo --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:2181

### Verify that was created
docker-compose exec kafka  \
  kafka-topics --describe --topic foo --zookeeper zookeeper:2181
  
### To publish some messages
docker-compose exec kafka  \
  bash -c "seq 10 | kafka-console-producer --request-required-acks 1 --broker-list localhost:29092 --topic foo && echo 'Produced 10 messages.'"
  
### Read messages using kafka console consumer
docker-compose exec kafka  \
  kafka-console-consumer --bootstrap-server localhost:29092 --topic foo  --from-beginning --max-messages 42
