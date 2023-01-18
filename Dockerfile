FROM golang:1.19.4-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

EXPOSE 8080

RUN apt update && apt install -y postgresql 
RUN make setup build

CMD ["make", "migrateup", "run"]
