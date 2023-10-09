FROM golang:1.21.2

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

ENV GO_ENV=dev

CMD ["air", "-c", ".air.toml"]