LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-note-api

generate-note-api:
	mkdir -p pkg/auth
	protoc --proto_path api/auth_proto \
	--go_out=pkg/auth --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/auth_proto/auth.proto

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go

copy-to-server:
	scp service_linux root@185.91.53.71:

docker-build-and-push:
	sudo docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/rxc/auth-serv .
	sudo docker login -u token -p CRgAAAAA8eHjytZYbhvTuUaSQpXlPEIZAIA9QY7k cr.selcloud.ru/rxc
	sudo docker push cr.selcloud.ru/rxc/auth-serv