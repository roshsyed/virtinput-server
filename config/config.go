package config

type Config struct {
	Host          string
	Port          string
	DevicePath    string
	VKeyboardName string
	VMouseName    string
}

func NewDefaultConfig() Config {
	return Config{
		Host:          "0.0.0.0",
		Port:          "50051",
		DevicePath:    "/dev/uinput",
		VKeyboardName: "virtinput-server Virtual Keyboard",
		VMouseName:    "virtinput-server Virtual Mouse",
	}
}
