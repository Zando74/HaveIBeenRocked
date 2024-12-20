package mock

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/grpc_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/test/bufconn"
)

type GrpcServerMock struct {
	grpc_proto.UnimplementedRawPasswordListUploadServer
	DB map[string]bool
}

func (s *GrpcServerMock) UploadRawPasswordList(stream grpc_proto.RawPasswordListUpload_UploadRawPasswordListServer) error {
	for {
		passwords, err := stream.Recv()
		if err != nil {

			if err == io.EOF {
				return stream.SendAndClose(&grpc_proto.Status{Success: true, Message: "Upload complete"})
			}

			return stream.SendAndClose(&grpc_proto.Status{Success: false, Message: err.Error()})
		}

		for _, password := range passwords.GetPasswords() {
			s.DB[string(password)] = true
		}
	}

}

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

func InitTestGrpcServerConnection(mockServer *GrpcServerMock) (*grpc.ClientConn, func()) {
	lis := bufconn.Listen(bufSize)
	resolver.Register(&bufconnResolver{lis: lis})

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	grpc_proto.RegisterRawPasswordListUploadServer(grpcServer, mockServer)

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
