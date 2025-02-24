FROM golang:1.23.0
RUN apt-get update && apt-get install -y gcc libc-dev sqlite3 libsqlite3-dev
WORKDIR /app
COPY . .
RUN go mod vendor 
RUN go test ./...
RUN go build -o twitter ./cmd/twitter/main.go
EXPOSE 8080
CMD ["./twitter"]
