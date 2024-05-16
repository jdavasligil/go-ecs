package ecs_test

import (
	"testing"

	"github.com/jdavasligil/go-ecs"
	"github.com/jdavasligil/go-ecs/pkg/testutil"
)

func TestEcs(t *testing.T) {
	world := ecs.NewWorld(ecs.WorldOptions{
		EntityLimit:    1024,
		RecycleLimit:   1024,
		ComponentLimit: 255,
	})
	player := world.NewEntity()
	t.Run("Initialize", func(t *testing.T) {
		ecs.Initialize[Position](&world)
		ecs.Initialize[Velocity](&world)
		ecs.Initialize[Health](&world)
		ecs.Initialize[CombatTag](&world)
	})
	t.Run("AddRemove", func(t *testing.T) {
		testutil.AssertEqual(t, ecs.Add(&world, player, Position{1.0, 2.0, 3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Velocity{-1.0, -2.0, -3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Health{16}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, CombatTag{}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, CombatTag{}), false)
		testutil.AssertEqual(t, ecs.Add(&world, player, DeadTag{}), false)

		testutil.AssertEqual(t, ecs.Remove[Position](&world, player), true)
		testutil.AssertEqual(t, ecs.Remove[Velocity](&world, player), true)
		testutil.AssertEqual(t, ecs.Remove[Health](&world, player), true)
		testutil.AssertEqual(t, ecs.Remove[CombatTag](&world, player), true)
		testutil.AssertEqual(t, ecs.Remove[DeadTag](&world, player), false)
		testutil.AssertEqual(t, ecs.Remove[Health](&world, player), false)
	})
}
