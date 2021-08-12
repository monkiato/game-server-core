package net

import (
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"net"
)

type GrpcModule struct {
	Server *grpc.Server
	port   int
}

func NewGrpcModule(port int) *GrpcModule {
	encoding.RegisterCodec(flatbuffers.FlatbuffersCodec{})

	var opts []grpc.ServerOption
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
