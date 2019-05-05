package version_test

import (
	"testing"

	. "github.com/goldeneggg/lsec2/cmd/lsec2/version"
)

func TestVersion(t *testing.T) {
	exp := "0.2.8"
	if VERSION != "0.2.8" {
		t.Errorf("expected: %#v, but actual: %#v", exp, VERSION)
	}
}
