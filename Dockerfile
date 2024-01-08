
FROM golang:alpine
WORKDIR /
COPY . .
RUN go mod download
RUN go build -o monitoring-service main.go
CMD ["./monitoring-service"]

