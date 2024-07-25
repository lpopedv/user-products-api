FROM golang:1.22-alpine

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

CMD ["air", "-c", ".air.toml"]
