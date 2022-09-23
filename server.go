package main

import (
	"context"
	"fmt"
	"net"

	hellov1 "github.com/echovl/hello-grpc/gen/proto/hello/v1"
	"google.golang.org/grpc"
)

type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{addr}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.Addr, err)
	}

	server := grpc.NewServer()
	hellov1.RegisterHelloServiceServer(server, &helloServiceServer{})
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

type helloServiceServer struct {
	hellov1.UnimplementedHelloServiceServer
}

func (s *helloServiceServer) Hello(ctx context.Context, req *hellov1.HelloRequest) (*hellov1.HelloResponse, error) {
	username := req.GetUsername()
	return &hellov1.HelloResponse{Msg: fmt.Sprintf("Hello %s", username)}, nil
}
