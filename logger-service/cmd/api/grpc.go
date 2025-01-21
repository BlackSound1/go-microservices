package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/BlackSound1/go-microservices/logger/data"
	"github.com/BlackSound1/go-microservices/logger/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

// WriteLog writes the given log entry to the database
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {

	input := req.GetLogEntry()

	// Write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// Return the response
	res := &logs.LogResponse{Result: "logged"}
	return res, nil
}

// grpcListen starts the gRPC server and listens on the configured GRPC_PORT (default 50001).
// It registers the LogServer with the gRPC server and logs a message when the server is
// started
func (app *Config) grpcListen() {

	// Create a listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", GRPC_PORT))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Create a gRPC server
	s := grpc.NewServer()

	// Register the LogServer
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", GRPC_PORT)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
