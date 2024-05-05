FROM golang:alpine AS builder
WORKDIR /app
ADD go.mod .
COPY . .
RUN go build -o bot .

FROM ubuntu
WORKDIR /app
COPY --from=builder /app/bot /app/bot
# copy the ca-certificate.crt from the build stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./bot"]

