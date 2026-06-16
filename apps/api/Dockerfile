FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum* ./
RUN go mod download
COPY . .

RUN swag init \
    --generalInfo main.go \
    --output ./docs \
    --parseDependency \
    --parseInternal

RUN go build -o server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]