# builder image
FROM golang:1.13-alpine3.11 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o . server/main.go


# generate clean, final image for end users
FROM alpine:3.11.3
COPY --from=builder /build/main ./server
ADD ./config /config/

EXPOSE 3000

# executable
ENTRYPOINT [ "./server" ]