package main

import (
	"context"
	"time"

	pb "github.com/eriktate/spinup"
)

// EchoService implements the EchoServer interface in echo.pb.go
type EchoService struct{}

// Echo is the only function defined on the EchoServer interface.
func (e *EchoService) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Msg:      in.GetMsg(),
		UnixTime: time.Now().Unix(),
	}, nil
}
