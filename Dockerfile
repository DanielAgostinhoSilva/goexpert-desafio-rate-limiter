FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o desafio-rate-limiter .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/.env .
COPY --from=builder /app/desafio-rate-limiter .
CMD ["./desafio-rate-limiter"]