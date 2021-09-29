
topic:
	docker exec -it redpanda_redpanda_1  rpk topic create topic-A --replicas 1

list:
	docker exec -it redpanda_redpanda_1  rpk topic list

consume:
	docker exec -it redpanda_redpanda_1 rpk topic consume topic-A

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
	docker exec -it redpanda_chameleon_1 \
		./client redpanda:9092