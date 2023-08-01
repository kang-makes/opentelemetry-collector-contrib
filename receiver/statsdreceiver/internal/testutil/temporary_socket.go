package testutil

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
