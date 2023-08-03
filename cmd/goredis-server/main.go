package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/chensylz/goredis/internal/server"
)

func main() {
	address := "localhost:6379"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	log.Printf("Server is listening on %s\n", address)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		log.Println("Received termination signal. Shutting down gracefully...")
		cancel()
		listener.Close()
	}()

	acceptConnections(ctx, listener)
}

func acceptConnections(ctx context.Context, listener net.Listener) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				continue
			}
			go handleConnection(ctx, conn)
		}
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	select {
	case <-ctx.Done():
		return
	default:
		server.HandleConnection(conn)
	}
}
