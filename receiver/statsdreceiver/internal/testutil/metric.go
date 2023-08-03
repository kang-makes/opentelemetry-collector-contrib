// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package testutil // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver/internal/transport/testutil"

import (
	"fmt"
)

// Metric contains the metric fields for a StatsDTestClient message.
type Metric struct {
	Name  string
	Value string
	Type  string
}

// String formats a Metric into a StatsDTestClient message.
func (m Metric) String() string {
	return fmt.Sprintf("%s:%s|%s", m.Name, m.Value, m.Type)
}
