FROM golang:1.19.4-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

EXPOSE 8080

RUN apk update && apk add postgresql-client make
RUN make setup build

CMD ["make", "migrateup", "run"]
