# build the application from an image
FROM golang:1.21.0-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o authApp ./cmd/api

RUN chmod +x /app/authApp

# build application from a scratch image
# this enables use to have a very slim build at the end of the day

FROM alpine:3.19

RUN mkdir /app

COPY --from=builder /app/authApp /app

CMD ["/app/authApp"]