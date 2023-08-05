package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chensylz/goredis/config"
	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/server/handler"
)

type App struct {
	config.Config
	handler handler.Handler
}

func New(config config.Config, handler handler.Handler) *App {
	return &App{
		Config:  config,
		handler: handler,
	}
}

func (a *App) Run() {
	listener, err := net.Listen("tcp", a.Address())
	if err != nil {
		logger.Infof("Error starting server: %v", err)
		return
	}
	defer listener.Close()

	logger.Infof("Server is listening on %s\n", a.Address())

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
	a.acceptConnections(ctx, listener)
}

func (a *App) Address() string {
	return fmt.Sprintf("%s:%d", a.Bind, a.Port)
}

func (a *App) acceptConnections(ctx context.Context, listener net.Listener) {
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
				a.handleConnection(ctx, conn, &wg)
			}()
		}
	}
}

func (a *App) handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer func() {
		_ = conn.Close()
		_ = wg.Done
	}()

	select {
	case <-ctx.Done():
		return
	default:
		a.handler.Handle(ctx, conn)
	}
}
