FROM golang:alpine AS builder
WORKDIR /app
ADD go.mod .
COPY . .
RUN go build -o bot .

FROM ubuntu
WORKDIR /app
COPY --from=builder /app/bot /app/bot
CMD ["./bot"]

