FROM golang:1.23.2-alpine AS builder

COPY . /github.com/IRXCI/auth/source/
WORKDIR /github.com/IRXCI/auth/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/IRXCI/auth/source/bin/crud_server .

CMD ["./crud_server"]