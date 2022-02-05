package rpc

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
	identOnce   sync.Once
	identConn   *grpc.ClientConn
	identLookup rpc.LookupServiceClient
	identExit   func() error
)

func LookupService() rpc.LookupServiceClient {
	identOnce.Do(func() {
		addr := net.JoinHostPort(config.IdentityHost(), config.IdentityPort())

		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		identConn = conn
		identLookup = pb.NewLookupServiceClient(identConn)
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
