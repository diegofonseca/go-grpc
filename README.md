# GRPC

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=. --go-grpc_out=. proto/course_category.proto

go run cmd/grpcServer/main.go

evans -r repl