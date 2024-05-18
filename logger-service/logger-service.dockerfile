# build the application from an image
FROM golang:1.21.0-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o loggerApp ./cmd/api

RUN chmod +x /app/loggerApp

# build application from a scratch image
# this enables use to have a very slim build at the end of the day

FROM alpine:3.19

RUN mkdir /app

COPY --from=builder /app/loggerApp /app

CMD ["/app/loggerApp"]