FROM arm64v8/golang:1.23

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

EXPOSE 8080

RUN go build -v -o /usr/local/bin/app ./cmd/test_server

EXPOSE 8080

RUN chmod +x /usr/local/bin/app

CMD ["/usr/local/bin/app"]
