Installing go so I could build it involved these steps

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export GO_PATH=~/go
export PATH=$PATH:/$GO_PATH/bin
protoc -I=./protos --go_out=./protos ./protos/llm_request.proto
```
