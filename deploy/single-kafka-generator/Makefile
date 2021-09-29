
topic:
	docker exec -it single-kafka-generator_kafka_1 /opt/bitnami/kafka/bin/kafka-topics.sh \
		--create --topic topic-A --bootstrap-server localhost:9092

list:
	docker exec -it single-kafka-generator_kafka_1 /opt/bitnami/kafka/bin/kafka-topics.sh \
		--list --bootstrap-server localhost:9092

consume:
	docker exec -it single-kafka-generator_kafka_1 /opt/bitnami/kafka/bin/kafka-console-consumer.sh \
		--topic topic-A --bootstrap-server localhost:9092

generate:
	curl -X 'POST' \
  		'http://localhost:3000/generators/people/start' \
  		-H 'accept: application/json' \
  		-d ''

stop:
	curl -X 'POST' \
		'http://localhost:3000/generators/people/stop' \
		-H 'accept: application/json' \
		-d ''

test:
	docker exec -it single-kafka-generator_chameleon_1 \
		./client kafka:9092
 