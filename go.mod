module github.com/tangxusc/crypt

go 1.15

require (
	cloud.google.com/go/firestore v1.1.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/hashicorp/consul/api v1.1.0
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/exp v0.0.0-20200331195152-e8c3332aa8e5 // indirect
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd // indirect
	google.golang.org/api v0.13.0
	google.golang.org/grpc v1.32.0
)

replace google.golang.org/grpc v1.32.0 => google.golang.org/grpc v1.26.0
