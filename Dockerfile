FROM golang:1.23-alpine

WORKDIR ./app

EXPOSE 8080
CMD ["./bin/api"]

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/api ./cmd/server \
    && go build -o ./bin/migrate ./cmd/migration/ \
    && go build -o ./bin/seed ./cmd/seed/ \