FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o golang-kafka-kubernetes-microservice cmd/server/main.go

EXPOSE 8080

CMD ["./golang-kafka-kubernetes-microservice"]