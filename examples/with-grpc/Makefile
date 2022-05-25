generate:
	 protoc -I=hatmaker/ --go_opt=module=github.com/supertokens/supertokens-golang/examples/with-grpc --go_out=hatmaker --go-grpc_opt=module=github.com/supertokens/supertokens-golang/examples/with-grpc --go-grpc_out=hatmaker hatmaker/*.proto

build:
	go build -o bin/server ./server
	go build -o bin/client ./client