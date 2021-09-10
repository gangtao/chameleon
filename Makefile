
PWD=$(shell pwd)

.PHONY: gen

all: gen run 

run:
	go run server/main.go

gen:
	$(HOME)/go/bin/swag init -g server/main.go 