package transport // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver/internal/transport"

// Transport is a set of constants of the transport supported by this receiver.
type Transport string

const (
	UDP  Transport = "udp"
	UDP4 Transport = "udp4"
	UDP6 Transport = "udp6"
)

// Create a Transport based on the transport string or return error if it is not supported.
func NewTransport(ts string) Transport {
	trans := Transport(ts)
	switch trans {
	case UDP, UDP4, UDP6:
		return trans
	}
	return Transport("")
}

// Returns the string of this transport.
func (trans Transport) String() string {
	switch trans {
	case UDP, UDP4, UDP6:
		return string(trans)
	}
	return ""
}

// Returns true if the transport us packet based.
func (trans Transport) IsPacketTransport() bool {
	switch trans {
	case UDP, UDP4, UDP6:
		return true
	}
	return false
}
