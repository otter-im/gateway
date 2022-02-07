package service

import (
	"github.com/otter-im/auth/internal/config"
	"github.com/otter-im/identity/pkg/rpc"
	pb "github.com/otter-im/identity/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
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
		addr := net.JoinHostPort(config.IdentityHost(), config.IdentityPort())

		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
