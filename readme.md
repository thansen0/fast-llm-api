# Setup

Installing go so I could build it involved these steps. Also [a good tutorial](https://protobuf.dev/getting-started/gotutorial/).

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export GO_PATH=~/go
export PATH=$PATH:/$GO_PATH/bin
protoc -I=./protos --go_out=./protos ./protos/llm_request.proto
```

I had more luck with 

```
protoc -I. --go-grpc_out=. llm_request.proto
```

for some unknown reason, must test later whether this was because I set my export path variables or because it's better somehow
