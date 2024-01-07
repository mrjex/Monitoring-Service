FROM golang:alpine
WORKDIR /
COPY . .
RUN go mod download
RUN go build -o monitor-service main.go
CMD ["./monitor-service"]

