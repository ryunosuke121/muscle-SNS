FROM golang:1.21.2

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /build/app ./cmd/main.go

EXPOSE 80

CMD ["/build/app"]