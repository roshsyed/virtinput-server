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

type MouseServer struct {
	vdevicev1.UnimplementedMouseServiceServer
	vm uinput.Mouse
	mu sync.Mutex
}

func (s *MouseServer) UseMouse(ctx context.Context, req *vdevicev1.MouseEvent) (*emptypb.Empty, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var err error

	switch req.Action {
	case vdevicev1.MouseAction_MOUSE_ACTION_UNSPECIFIED:
		return nil, status.Errorf(codes.InvalidArgument, "mouse action unspecified")

	// --- Directional moves (scalar pixels) ---
	case vdevicev1.MouseAction_MOUSE_MOVE_LEFT:
		err = s.vm.MoveLeft(req.Pixel)
	case vdevicev1.MouseAction_MOUSE_MOVE_RIGHT:
		err = s.vm.MoveRight(req.Pixel)
	case vdevicev1.MouseAction_MOUSE_MOVE_UP:
		err = s.vm.MoveUp(req.Pixel)
	case vdevicev1.MouseAction_MOUSE_MOVE_DOWN:
		err = s.vm.MoveDown(req.Pixel)
	// --- Vector move (dx, dy) ---
	case vdevicev1.MouseAction_MOUSE_MOVE:
		err = s.vm.Move(req.X, req.Y)
	// --- Clicks ---
	case vdevicev1.MouseAction_MOUSE_LEFT_CLICK:
		err = s.vm.LeftClick()
	case vdevicev1.MouseAction_MOUSE_RIGHT_CLICK:
		err = s.vm.RightClick()
	case vdevicev1.MouseAction_MOUSE_MIDDLE_CLICK:
		err = s.vm.MiddleClick()
	// --- Button press/release ---
	case vdevicev1.MouseAction_MOUSE_LEFT_PRESS:
		err = s.vm.LeftPress()
	case vdevicev1.MouseAction_MOUSE_LEFT_RELEASE:
		err = s.vm.LeftRelease()
	case vdevicev1.MouseAction_MOUSE_RIGHT_PRESS:
		err = s.vm.RightPress()
	case vdevicev1.MouseAction_MOUSE_RIGHT_RELEASE:
		err = s.vm.RightRelease()
	case vdevicev1.MouseAction_MOUSE_MIDDLE_PRESS:
		err = s.vm.MiddlePress()
	case vdevicev1.MouseAction_MOUSE_MIDDLE_RELEASE:
		err = s.vm.MiddleRelease()
	// --- Wheel ---
	case vdevicev1.MouseAction_MOUSE_WHEEL:
		err = s.vm.Wheel(req.Horizontal, req.Delta)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported mouse action: %v", req.Action)
	}
	if err != nil {
		// Preserve the underlying error, but return a stable gRPC status.
		return nil, status.Errorf(codes.Internal, "%s failed: %v", req.Action, err)
	}
	return &emptypb.Empty{}, nil
}

func NewMouseServer(vm uinput.Mouse) *MouseServer {
	return &MouseServer{vm: vm}
}
