
BIN_NAME ?= chameleon
VERSION ?= 0.1
IMAGE_NAME ?= $(BIN_NAME):$(VERSION)
DOCKER_ID_USER ?= naughtytao

DATE=$(shell date '+%Y%m%d')
FULLNAME=$(DOCKER_ID_USER)/${IMAGE_NAME}.${DATE}

PWD=$(shell pwd)

.PHONY: gen docker

all: gen run 

run:
	go run server/main.go

build:
	go build server/main.go

gen:
	$(HOME)/go/bin/swag init -g server/main.go 

test:
	go test ./generator/
	go test ./source/
	go test ./sink/
	go test ./handlers/

docker: Dockerfile
	docker build -t $(IMAGE_NAME) .

docker_run:
	docker run -p 3000:3000 $(IMAGE_NAME)

push:
	docker tag $(IMAGE_NAME) $(FULLNAME)
	docker push $(FULLNAME)

test_create:
	curl -X 'POST' \
		'http://localhost:3000/generators' \
		-H 'accept: application/json' \
		-H 'Content-Type: application/json' \
		-d '{
			"name": "testconfig3",
			"sink": {
				"name": "sinkname",
				"type": "kafka",
				"config": {
				"brokers": [
					"kafka1:9092",
					"kafka2:9093",
					"kafka3:9094"
				],
				"topic": "topic-A"
				}
			},
			"source": {
				"name": "sourcename",
				"timestamp_field": "f1",
				"batch_size": 2,
				"concurrency": 2,
				"interval": [
				10000,
				20000
				],
				"fields": [
					{
						"name": "f1",
						"type": "timestamp"
					},
					{
						"name": "f2",
						"type": "string",
						"range": [
						"a",
						"b",
						"c"
						]
					},
					{
						"name": "f3",
						"type": "int",
						"range": [
						1,
						5,
						30
						]
					},
					{
						"name": "f4",
						"type": "float",
						"range": [
						1.2,
						5.5,
						30.2
						]
					}
				]
			}
		}'