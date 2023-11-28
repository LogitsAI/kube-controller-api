package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/LogitsAI/kube-controller-api/controllerpb"
	"github.com/LogitsAI/kube-controller-api/manager"
	"github.com/go-logr/logr/slogr"
	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func main() {
	// Configure the default slog logger.
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))
	// Set controller-runtime to log to slog with the same settings.
	log.SetLogger(slogr.NewLogr(handler))

	kubeConfig := config.GetConfigOrDie()
	cmServer := manager.NewControllerManagerServer(kubeConfig)

	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(fmt.Errorf("can't listen on gRPC port: %w", err))
	}

	grpcServer := grpc.NewServer()
	controllerpb.RegisterControllerManagerServer(grpcServer, cmServer)
	go func() {
		slog.Info("Starting gRPC server", "address", listener.Addr())
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()
	defer grpcServer.Stop()

	// Wait for the process to be terminated.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan

	slog.Info("Received signal; shutting down", "signal", sig)
}
