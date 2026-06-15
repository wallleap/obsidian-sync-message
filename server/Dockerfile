FROM docker.1ms.run/golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn,direct go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/main.go

FROM docker.1ms.run/alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

COPY --from=builder /app/server ./server

RUN mkdir -p logs data/uploads

ENV SERVER_PORT=8080

EXPOSE 8080

CMD ["./server"]
