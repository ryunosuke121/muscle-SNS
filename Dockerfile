FROM golang:1.21.2

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/rubenv/sql-migrate/...@latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# マイグレーションスクリプトの実行権限付与
RUN chmod +x database/scripts/up.sh
RUN chmod +x database/scripts/down.sh

ENV GO_ENV=dev

CMD ["air", "-c", ".air.toml"]