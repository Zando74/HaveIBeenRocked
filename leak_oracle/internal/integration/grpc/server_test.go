package grpc_test

import (
	"context"
	"log"
	"net"

	grpc_controller "github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/controller/grpc"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/grpc_proto"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/test/bufconn"
)

var (
	LeakedHashRepositoryImpl = repository.NewHashLeakedRepositoryImpl()
)

const bufSize = 101024 * 1024
const bufnetScheme = "bufnet"

type bufconnResolver struct {
	lis *bufconn.Listener
}

func (r *bufconnResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	go func() {
		state := resolver.State{
			Addresses: []resolver.Address{{Addr: target.Endpoint()}},
		}
		cc.UpdateState(state)
	}()
	return r, nil
}

func (r *bufconnResolver) Scheme() string                        { return bufnetScheme }
func (r *bufconnResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *bufconnResolver) Close()                                {}

func InitTestServer() (*grpc.ClientConn, func()) {
	lis := bufconn.Listen(bufSize)
	resolver.Register(&bufconnResolver{lis: lis})

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	grpc_proto.RegisterRawPasswordListUploadServer(grpcServer, &grpc_controller.PasswordListUploadServer{
		ProcessPasswordBatchUseCase: &usecase.ProcessPasswordBatch{
			HashPasswordBatchUseCase: &usecase.HashPasswordBatch{
				LeakedHashRepository: LeakedHashRepositoryImpl,
			},
			HashLeakedRepository: LeakedHashRepositoryImpl,
		},
	})

	go func(g *grpc.Server) {

		if err := g.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}(grpcServer)

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient(bufnetScheme+":///test",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(),
	)

	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		grpcServer.Stop()
	}

	return conn, closer
}
