package service

import (
	"github.com/otter-im/gateway/internal/config"
	"github.com/otter-im/identity/pkg/rpc"
	pb "github.com/otter-im/identity/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var (
	identConnOnce   sync.Once
	identLookupOnce sync.Once
	identConn       *grpc.ClientConn
	identLookup     rpc.LookupServiceClient
	identExit       func() error
)

func IdentityConn() *grpc.ClientConn {
	identConnOnce.Do(func() {
		conn, err := grpc.Dial(config.IdentityAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		conn.Connect()
		identConn = conn
	})
	return identConn
}

func LookupService() rpc.LookupServiceClient {
	identLookupOnce.Do(func() {
		identLookup = pb.NewLookupServiceClient(IdentityConn())
	})
	return identLookup
}

func ExitHook() error {
	if identConn != nil {
		if err := identConn.Close(); err != nil {
			return err
		}
	}
	return nil
}
