package fast_test

import (
	"slices"
	"testing"

	"github.com/nilathedragon/spamscale/pkg/thirdparty/fast"
)

func TestIsBlocked(t *testing.T) {
	blocklist, err := fast.GetBlocklist()
	if err != nil {
		t.Fatalf("Failed to get blocklist: %v", err)
	}
	if len(blocklist) == 0 {
		t.Fatalf("Blocklist is empty")
	}
	if !slices.Contains(blocklist, "6549980962") {
		t.Fatalf("Known blocked user 6549980962 is not on the blocklist")
	}
}
