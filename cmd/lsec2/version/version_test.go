package version_test

import (
	"testing"

	. "github.com/goldeneggg/lsec2/cmd/lsec2/version"
)

func TestVersion(t *testing.T) {
	exp := "0.2.11"

	if VERSION != exp {
		t.Errorf("expected: %#v, but actual: %#v", exp, VERSION)
	}
}
