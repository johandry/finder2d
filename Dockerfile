# Builder image
FROM golang AS builder

WORKDIR /workspace/finder2d

ENV     GO111MODULE=on
COPY    go.mod .
COPY    go.sum .
RUN     go mod tidy && go mod download

COPY    . .
RUN     CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /finder2d cmd/main.go

# Application image
FROM alpine:3.9 AS application

COPY --from=builder /finder2d /app/

ENTRYPOINT [ "/app/finder2d" ]
