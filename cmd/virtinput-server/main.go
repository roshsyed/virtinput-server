package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	vdevicev1 "github.com/roshsyed/virtinput-server/api/vdevice/v1"
	vdev "github.com/roshsyed/virtinput-server/internal/server"

	"github.com/bendahl/uinput"
	"github.com/roshsyed/virtinput-server/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Default config from config/config.go
	conf := config.NewDefaultConfig()
	// TCP listener default 0.0.0.0:50051
	lis, err := net.Listen("tcp", net.JoinHostPort(conf.Host, conf.Port))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	// Create the virtual keyboard
	// The user running this program must have proper permissions on the device path
	keyboard, err := uinput.CreateKeyboard(conf.DevicePath, []byte(conf.VKeyboardName))
	if err != nil {
		log.Fatalf("vkeyboard: %v", err)
	}
	defer func() {
		if err := keyboard.Close(); err != nil {
			log.Printf("keyboard close failed: %v", err)
		}
	}()

	mouse, err := uinput.CreateMouse(conf.DevicePath, []byte(conf.VMouseName))
	if err != nil {
		log.Fatalf("vmouse: %v", err)
	}
	defer func() {
		if err := mouse.Close(); err != nil {
			log.Printf("mouse close failed: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	vdevicev1.RegisterKeyboardServiceServer(grpcServer, vdev.NewKeyboardServer(keyboard))
	vdevicev1.RegisterMouseServiceServer(grpcServer, vdev.NewMouseServer(mouse))

	log.Printf("gRPC server listening on %s:%s", conf.Host, conf.Port)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		log.Println("shutting down virtinput-server")
		grpcServer.GracefulStop()
		if err := lis.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("listener close error: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, net.ErrClosed) {
		log.Fatalf("serve: %v", err)
	}
}
