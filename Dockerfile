# builder image
FROM golang:1.17.1-alpine3.13 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./mains server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./mainc client/main.go


# generate clean, final image for end users
FROM alpine:3.13.6
COPY --from=builder /build/mains ./server
COPY --from=builder /build/mainc ./client
ADD ./config /config/

EXPOSE 3000

# executable
ENTRYPOINT [ "./server" ]