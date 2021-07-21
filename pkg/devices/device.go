package devices

type Device interface {
	GetConnectionInfo() ConnectionInfo
	GetTransceiverInfo() (*DeviceOpticTransceiver, error)
}

type ConnectionInfo struct {
	Hostname string
	Username string
	Password string
	Ciphers  []string
}

type DeviceOpticTransceiver struct {
	LinkStatus bool
	RXRSSI     float64
	TXRSSI     float64
}
