// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package testutil // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver/transport/testutil"

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func CreateTemporarySocket(t testing.TB, _ string) string {
	b := make([]byte, 10)
	rand.Read(b)
	return fmt.Sprintf("%s/%s", t.TempDir(), fmt.Sprintf("%x", b))
}
