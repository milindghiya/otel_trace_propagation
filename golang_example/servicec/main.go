package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/milindghiya/otel_trace_propagation/golang_example/hello"
	"github.com/milindghiya/otel_trace_propagation/golang_example/otel_utils"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	hello.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	fmt.Println(span.SpanContext().TraceID())
	return &hello.HelloResponse{Message: "\nHello from service C  TraceID: " + span.SpanContext().TraceID().String() + " SpanID: " + span.SpanContext().SpanID().String()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	om := otel_utils.InitOtelManager("ServiceC", "v0.0.0")
	otelShutdown, err := om.SetupOTelSDK(context.Background())
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	hello.RegisterHelloServiceServer(s, &server{})

	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
