// Package server implements the gRPC handlers for virtinput-server.
// It translates API requests into serialized uinput operations.
package server

import (
	"context"
	"sync"

	vdevicev1 "github.com/roshsyed/virtinput-server/api/vdevice/v1"

	"github.com/bendahl/uinput"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type KeyboardServer struct {
	vdevicev1.UnimplementedKeyboardServiceServer
	vk uinput.Keyboard
	mu sync.Mutex
}

func (s *KeyboardServer) SendKey(ctx context.Context, req *vdevicev1.KeyEvent) (*emptypb.Empty, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var err error

	switch req.Action {
	case vdevicev1.KeyAction_KEY_PRESS:
		err = s.vk.KeyPress(int(req.KeyCode))
	case vdevicev1.KeyAction_KEY_DOWN:
		err = s.vk.KeyDown(int(req.KeyCode))
	case vdevicev1.KeyAction_KEY_UP:
		err = s.vk.KeyUp(int(req.KeyCode))
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported keyboard action: %v", req.Action)
	}
	if err != nil {
		// Preserve the underlying error, but return a stable gRPC status.
		return nil, status.Errorf(codes.Internal, "%s failed: %v", req.Action, err)
	}
	return &emptypb.Empty{}, nil
}

func NewKeyboardServer(vk uinput.Keyboard) *KeyboardServer {
	return &KeyboardServer{vk: vk}
}
