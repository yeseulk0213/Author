package server

import (
	"net"

	"github.com/golang/glog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"gitlab.com/promptech1/infuser-author/app/ctx"
	"gitlab.com/promptech1/infuser-author/handler"
	grpc_author "gitlab.com/promptech1/infuser-author/infuser-protobuf/gen/proto/author"
	"google.golang.org/grpc"
)

// Server is an main application object that shared (read-only) to application modules
type Server struct {
	ctx        *ctx.Context
	grpcServer *grpc.Server
}

// New constructor
func New(c *ctx.Context) *Server {
	s := new(Server)
	s.ctx = c
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	return s
}

func (s *Server) Run(network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	appTokenHandler := handler.NewAppTokenHandler(s.ctx)

	// Token 기반의 인증 처리
	grpc_author.RegisterApiAuthServiceServer(s.grpcServer, newApiAuthServer(appTokenHandler))

	go func() {
		defer s.grpcServer.GracefulStop()
		<-s.ctx.Context.Done()
	}()

	glog.Infof("start gRPC grpc at %s", address)
	return s.grpcServer.Serve(l)

}
