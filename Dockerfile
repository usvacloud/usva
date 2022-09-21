FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN make setup build migratesetup
CMD ["make", "migrateup", "run"]