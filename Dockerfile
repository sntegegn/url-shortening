FROM golang:latest

WORKDIR /app

COPY . . 

RUN go mod download

RUN make build

CMD ["./shortner"]