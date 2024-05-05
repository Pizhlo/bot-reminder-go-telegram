FROM golang:alpine AS builder
WORKDIR /app
ADD go.mod .
COPY . .
RUN go build -o bot .

FROM docker:dind
WORKDIR /app
COPY --from=builder /app/bot /app/bot
COPY .env .env
RUN apk add bash
COPY migration migration
COPY migrate.sh .
RUN chmod +x migrate.sh
CMD ["./bot"]

