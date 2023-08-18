
FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o receipt-points-calculator .

EXPOSE 4000

CMD ["./receipt-points-calculator"]