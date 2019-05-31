GOGOPATH=${GOPATH}/src
protoc --proto_path=${GOPATH}:${GOGOPATH}:./ --gofast_out=plugins=grpc:. *.proto