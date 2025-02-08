# build the application from an image
FROM golang:1.23.6-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# build application from a scratch image
# this enables use to have a very slim build at the end of the day

FROM alpine:3.21

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]