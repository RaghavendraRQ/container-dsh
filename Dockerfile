FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o container-dsh ./cmd/container-dsh

EXPOSE 8080

CMD ["./container-dsh"]