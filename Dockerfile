FROM golang:1.23.2-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

WORKDIR /app/src

RUN swag init
RUN go build -o message-automation ./main.go

EXPOSE 3030

CMD ["./message-automation"]
