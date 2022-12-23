FROM golang:1.19-alpine AS tools

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 DATABASE_URL=postgres://postgres:postgres@pgsql15/postgres?sslmode=disable

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD [ "go", "test", "-v", "-cover", "./..." ]