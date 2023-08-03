// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package transport // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver/internal/transport"

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/multierr"
)

type packetServer struct {
	packetConn net.PacketConn
	transport  Transport
	address    string
}

var (
	// Ensure that Server is implemented on UDP Server.
	_ (Server) = (*packetServer)(nil)

	ErrUnsupportedPacketTransport = errors.New("unsupported Packet transport")
)

// NewPacketServer creates a transport.Server using transports based on packets.
func NewPacketServer(transport Transport, address string) (Server, error) {
	if !transport.IsPacketTransport() {
		return nil, ErrUnsupportedPacketTransport
	}

	if transport.IsUnixTransport() {
		err := os.Remove(address)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("removing socket failed: %w", err)
		}
	}

	conn, err := net.ListenPacket(transport.String(), address)
	if err != nil {
		return nil, fmt.Errorf("starting to listen %s: %w", transport.String(), err)
	}

	return &packetServer{
		address:    address,
		packetConn: conn,
		transport:  transport,
	}, nil
}

// ListenAndServe starts the server ready to receive metrics.
func (psrv *packetServer) ListenAndServe(
	nextConsumer consumer.Metrics,
	reporter Reporter,
	transferChan chan<- Metric,
) error {
	if nextConsumer == nil || reporter == nil {
		return errNilListenAndServeParameters
	}

	buf := make([]byte, 65527) // max size for udp packet body (assuming ipv6)
	for {
		n, addr, err := psrv.packetConn.ReadFrom(buf)

		if psrv.transport.IsUnixTransport() && addr == nil {
			addr = &net.UnixAddr{
				Net:  "unixgram",
				Name: "UDS",
			}
		}

		if n > 0 {
			bufCopy := make([]byte, n)
			copy(bufCopy, buf)
			psrv.handlePacket(bufCopy, addr, transferChan)
		}
		if err != nil {
			reporter.OnDebugf("%s Transport (%s) - ReadFrom error: %v",
				psrv.transport,
				psrv.packetConn.LocalAddr(),
				err)
			var netErr net.Error
			if errors.As(err, &netErr) {
				if netErr.Timeout() {
					continue
				}
			}
			return err
		}
	}
}

// Close closes the server.
func (psrv *packetServer) Close() error {
	var errs error

	if psrv.transport.IsUnixTransport() {
		err := os.Remove(psrv.address)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = multierr.Append(errs, fmt.Errorf("removing socket failed: %w", err))
		}
	}

	err := psrv.packetConn.Close()
	if err != nil {
		errs = multierr.Append(errs, err)
	}

	return errs
}

// handlePacket is helper that parses the buffer and split it line by line to be parsed upstream.
func (psrv *packetServer) handlePacket(
	data []byte,
	addr net.Addr,
	transferChan chan<- Metric,
) {
	buf := bytes.NewBuffer(data)
	for {
		bytes, err := buf.ReadBytes((byte)('\n'))
		if errors.Is(err, io.EOF) {
			if len(bytes) == 0 {
				// Completed without errors.
				break
			}
		}
		line := strings.TrimSpace(string(bytes))
		if line != "" {
			transferChan <- Metric{line, addr}
		}
	}
}
