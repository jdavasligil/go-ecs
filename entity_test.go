package ecs_test

import (
	"testing"

	"github.com/jdavasligil/go-ecs"
)

func TestEntity(t *testing.T) {
	var id uint32 = (1 << 24) - 997
	var e ecs.Entity = ecs.NewEntity(id)

    // Test Creation
    if id != e.ID() {
        t.Errorf("Expected: %d, Got: %d", id, e.ID())
    }
    if uint8(0) != e.Version() {
        t.Errorf("Expected: %d, Got: %d", uint8(0), e.Version())
    }

    // Test Next
    e.Next()
    if id != e.ID() {
        t.Errorf("Expected: %d, Got: %d", id, e.ID())
    }
    if uint8(1) != e.Version() {
        t.Errorf("Expected: %d, Got: %d", uint8(1), e.Version())
    }

    // Test Rollover
    for i := 0; i < 255; i++ {
        e.Next()
    }
    if id != e.ID() {
        t.Errorf("Expected: %d, Got: %d", id, e.ID())
    }
    if uint8(0) != e.Version() {
        t.Errorf("Expected: %d, Got: %d", uint8(0), e.Version())
    }
}
