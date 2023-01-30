FROM golang:1.20rc1-alpine3.17

WORKDIR /usr/usermanager_api

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./bin/usermanager ./cmd/usermanager/

CMD [ "./bin/usermanager" ]