// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package testutil // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver/internal/testutil"

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// StatsDTestClient defines the properties of a StatsD connection.
type StatsDTestClient struct {
	transport string
	address   string
	conn      io.Writer
}

// NewStatsDTestClient creates a new StatsDTestClient instance to support the need for testing
// the statsdreceiver package and is not intended/tested to be used in production.
func NewStatsDTestClient(transport string, address string) (*StatsDTestClient, error) {
	statsd := &StatsDTestClient{
		transport: transport,
		address:   address,
	}

	err := statsd.connect()
	if err != nil {
		return nil, err
	}

	return statsd, nil
}

// connect populates the StatsDTestClient.conn
func (s *StatsDTestClient) connect() error {
	switch s.transport {
	case "udp":
		udpAddr, err := net.ResolveUDPAddr(s.transport, s.address)
		if err != nil {
			return err
		}
		s.conn, err = net.DialUDP(s.transport, nil, udpAddr)
		if err != nil {
			return err
		}
	case "unixgram":
		unixAddr, err := net.ResolveUnixAddr(s.transport, s.address)
		if err != nil {
			return err
		}
		s.conn, err = net.DialUnix(s.transport, nil, unixAddr)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown/unsupported transport: %s", s.transport)
	}

	return nil
}

// Disconnect closes the StatsDTestClient.conn.
func (s *StatsDTestClient) Disconnect() error {
	var err error
	if cl, ok := s.conn.(io.Closer); ok {
		err = cl.Close()
	}
	s.conn = nil
	return err
}

// SendMetric sends the input metric to the StatsDTestClient connection.
func (s *StatsDTestClient) SendMetric(metric Metric) error {
	_, err := io.Copy(s.conn, strings.NewReader(metric.String()))
	if err != nil {
		return fmt.Errorf("send metric on test client: %w", err)
	}
	return nil
}
