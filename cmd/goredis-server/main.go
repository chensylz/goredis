package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/server"
)

func main() {
	logger.SetupLog()
	address := "localhost:6379"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Infof("Error starting server: %v", err)
		return
	}
	defer listener.Close()

	logger.Infof("Server is listening on %s\n", address)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-signalCh
		logger.Infof("Received termination signal. Shutting down gracefully...")
		cancel()
		listener.Close()
	}()

	acceptConnections(ctx, listener)
}

func acceptConnections(ctx context.Context, listener net.Listener) {
	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go func() {
				wg.Add(1)
				handleConnection(ctx, conn, &wg)
			}()
		}
	}
}

func handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer func() {
		_ = conn.Close()
		_ = wg.Done
	}()

	select {
	case <-ctx.Done():
		return
	default:
		server.HandleConnection(conn)
	}
}
