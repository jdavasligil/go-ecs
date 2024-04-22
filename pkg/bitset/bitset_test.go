package bitset_test

import (
	"testing"

	"github.com/jdavasligil/go-ecs/pkg/bitset"
)

func TestNewUint64(t *testing.T) {
	bs := bitset.NewUint64(256)
	if len(bs) != 4 {
		t.Errorf("Expected len: 4, Got: %d\n", len(bs))
	}

	bs = bitset.NewUint64(255)
	if len(bs) != 4 {
		t.Errorf("Expected len: 4, Got: %d\n", len(bs))
	}

	bs = bitset.NewUint64(257)
	if len(bs) != 5 {
		t.Errorf("Expected len: 5, Got: %d\n", len(bs))
	}
}

func TestUint64(t *testing.T) {
	bs := bitset.NewUint64(256)
	bs.SetBit(5)
	if !bs.GetBit(5) {
		t.Errorf("Expected set bit at position 5. Got zero.")
	}

	bs.ClearBit(5)
	if bs.GetBit(5) {
		t.Errorf("Expected zero at position 5. Got one.")
	}
}
