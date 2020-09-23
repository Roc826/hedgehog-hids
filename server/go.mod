module hedgehog-hids-server

go 1.15

require (
	github.com/fatih/color v1.9.0
	github.com/smallnest/rpcx v0.0.0-20200917102714-42a82be8f8ab
	go.mongodb.org/mongo-driver v1.4.1
	google.golang.org/grpc/examples v0.0.0-20200921235902-d81def4352bc // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
