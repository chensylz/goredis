package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chensylz/goredis/config"
	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/server/handler"
	"github.com/chensylz/goredis/pkg/wait"
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
	var wg wait.Wait
	for {
		select {
		case <-ctx.Done():
			wg.WaitWithTimeout(10 * time.Second)
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go func() {
				wg.Add(1)
				a.handleConnection(ctx, conn)
				defer wg.Done()
			}()
		}
	}
}

func (a *App) handleConnection(ctx context.Context, conn net.Conn) {
	select {
	case <-ctx.Done():
		a.handler.Close()
		return
	default:
		a.handler.Handle(ctx, conn)
	}
}
