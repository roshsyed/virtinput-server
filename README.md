# virtinput-server
virtinput-server is a Linux daemon that exposes a gRPC API for creating and controlling virtual input devices (keyboard and mouse) using the Linux uinput subsystem.

It allows remote clients to send input events (keyboard, mouse movement, buttons, wheel) to a Linux host over the network.

*API is experimental and may change*

# Requirements
- Linux with uinput support
- Permission to create uinput devices (root or udev rule)

# Config
Default configs live in:
```
config/config.go
```

# Build
```
make generate
go build ./cmd/virtinput-server
```

# Run
```
sudo ./virtinput-server
```
Listens on :50051 by default.

# API
Protobuf definitions live in:
```
proto/vdevice/v1/vdevice.proto
```

Services:
- KeyboardService
- MouseService

Server reflection is enabled.

# Example
```
grpcurl -plaintext \
  -d '{"keyCode":30,"action":"KEY_PRESS"}' \
  localhost:50051 vdevice.v1.KeyboardService/SendKey
```

# License
MIT
