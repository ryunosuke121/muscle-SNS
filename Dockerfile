# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /src

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY ./ ./

# RUN chmod +x ./migrate.sh
# CMD [ "./migrate.sh" ]

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can (optionally) document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

FROM alpine:latest

WORKDIR /app

COPY --from=0 /src/app .

# Run
ENTRYPOINT ["./app"]