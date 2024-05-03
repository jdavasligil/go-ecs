package ecs_test

import (
	"testing"

	"github.com/jdavasligil/go-ecs"
	"github.com/jdavasligil/go-ecs/pkg/testutil"
)

type Velocity struct {
	x float32
	y float32
	z float32
}

func (c Velocity) ID() ecs.ComponentID {
	return 0
}

type Position struct {
	x float32
	y float32
	z float32
}

func (c Position) ID() ecs.ComponentID {
	return 1
}

func TestEcs(t *testing.T) {
	world := ecs.NewWorld()
	player := world.NewEntity()
	testutil.AssertEqual(t, ecs.Register(&world, player, Position{1.0, 2.0, 3.0}), true)
	testutil.AssertEqual(t, ecs.Register(&world, player, Velocity{-1.0, -2.0, -3.0}), true)
	pos, ok := ecs.Query[Position](&world, player)
	testutil.AssertEqual(t, ok, true)
	testutil.AssertEqual(t, pos.x, 1.0)
	vel, ok := ecs.Query[Velocity](&world, player)
	testutil.AssertEqual(t, ok, true)
	testutil.AssertEqual(t, vel.x, -1.0)
	mutVel, ok := ecs.MutQuery[Velocity](&world, player)
	testutil.AssertEqual(t, ok, true)
	mutVel.x = 0.0
	vel, ok = ecs.Query[Velocity](&world, player)
	testutil.AssertEqual(t, ok, true)
	testutil.AssertEqual(t, vel.x, 0.0)
}
