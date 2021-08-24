package net

import (
	"context"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"net"
)

type GrpcModule struct {
	Server *grpc.Server
	port   int
}

func NewGrpcModule(port int) *GrpcModule {
	encoding.RegisterCodec(flatbuffers.FlatbuffersCodec{})

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			logrus.Printf("intercepted call %s", info.FullMethod)
			metadata, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				logrus.Error("error trying to fetch metadata")
				return nil, fmt.Errorf("error trying to fetch metadata from request %s", info.FullMethod)
			}
			logrus.Printf("metadata length: %d", metadata.Len())
			for key, value := range metadata {
				logrus.Printf("metadata[%s] = %s", key, value)
			}
			return handler(ctx, req)
		}),
	}
	grpcServer := grpc.NewServer(opts...)

	return &GrpcModule {
		Server: grpcServer,
		port: port,
	}

}

func (m *GrpcModule) Run() {
	addr := fmt.Sprintf(":%d", m.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	logrus.Printf("running tcp server at %s", addr)
	m.Server.Serve(lis)
}
