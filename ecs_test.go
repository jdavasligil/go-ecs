package ecs_test

import (
	"testing"

	"github.com/jdavasligil/go-ecs"
	"github.com/jdavasligil/go-ecs/pkg/testutil"
)

const (
	VelocityID ecs.ComponentID = iota
	PositionID
	HealthID
)

type Velocity struct {
	x float32
	y float32
	z float32
}

func (c Velocity) ID() ecs.ComponentID {
	return VelocityID
}

type Position struct {
	x float32
	y float32
	z float32
}

func (c Position) ID() ecs.ComponentID {
	return PositionID
}

type Health struct {
	hp int
}

func (c Health) ID() ecs.ComponentID {
	return HealthID
}

func TestEcs(t *testing.T) {
	t.Run("Initialize", func(t *testing.T) {
		world := ecs.NewWorld()
		ecs.Initialize[Position](&world)
		ecs.Initialize[Velocity](&world)
		ecs.Initialize[Health](&world)
		player := world.NewEntity()
		testutil.AssertEqual(t, ecs.Add(&world, player, Position{1.0, 2.0, 3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Velocity{-1.0, -2.0, -3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Health{16}), true)
	})

	t.Run("RegisterQuery", func(t *testing.T) {
		world := ecs.NewWorld()
		player := world.NewEntity()
		testutil.AssertEqual(t, ecs.Register(&world, player, Position{1.0, 2.0, 3.0}), true)
		testutil.AssertEqual(t, ecs.Register(&world, player, Velocity{-1.0, -2.0, -3.0}), true)
		pos, ok := ecs.Get[Position](&world, player)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, pos.x, 1.0)
		vel, ok := ecs.Get[Velocity](&world, player)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, vel.x, -1.0)
		mutVel, ok := ecs.GetMut[Velocity](&world, player)
		testutil.AssertEqual(t, ok, true)
		mutVel.x = 0.0
		vel, ok = ecs.Get[Velocity](&world, player)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, vel.x, 0.0)
	})

	t.Run("AddRemove", func(t *testing.T) {
		world := ecs.NewWorld()
		ecs.Initialize[Position](&world)
		ecs.Initialize[Velocity](&world)
		ecs.Initialize[Health](&world)
		player := world.NewEntity()
		testutil.AssertEqual(t, ecs.Add(&world, player, Position{1.0, 2.0, 3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Velocity{-1.0, -2.0, -3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Health{16}), true)
		hp, ok := ecs.Get[Health](&world, player)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, hp.hp, 16)
		testutil.AssertEqual(t, ecs.Remove[Health](&world, player), player)
		testutil.AssertEqual(t, ecs.Remove[Health](&world, player), 0)
		hp, ok = ecs.Get[Health](&world, player)
		testutil.AssertEqual(t, ok, false)
	})

	t.Run("Query", func(t *testing.T) {
		world := ecs.NewWorld()
		ecs.Initialize[Position](&world)
		ecs.Initialize[Velocity](&world)
		ecs.Initialize[Health](&world)
		player := world.NewEntity()
		npc1 := world.NewEntity()
		npc2 := world.NewEntity()
		testutil.AssertEqual(t, ecs.Add(&world, player, Position{1.0, 2.0, 3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Velocity{-1.0, -2.0, -3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Health{16}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc1, Position{1.0, 1.0, 1.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc2, Position{2.0, 2.0, 2.0}), true)
		es, ps := ecs.Query[Position](&world)
		for i, e := range es {
			switch e {
			case player:
				testutil.AssertEqual(t, ps[i].z, 3)
			case npc1:
				testutil.AssertEqual(t, ps[i].z, 1)
			case npc2:
				testutil.AssertEqual(t, ps[i].z, 2)
			}
		}
	})

	t.Run("Query2", func(t *testing.T) {
		world := ecs.NewWorld()
		ecs.Initialize[Position](&world)
		ecs.Initialize[Velocity](&world)
		ecs.Initialize[Health](&world)
		player := world.NewEntity()
		npc1 := world.NewEntity()
		npc2 := world.NewEntity()
		testutil.AssertEqual(t, ecs.Add(&world, player, Position{1.0, 2.0, 3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Health{16}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc1, Position{1.0, 1.0, 1.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc1, Health{14}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc2, Position{2.0, 2.0, 2.0}), true)
		es := ecs.Query2[Position, Health](&world)
		testutil.AssertEqual(t, len(es), 2)
		for _, e := range es {
			switch e {
			case player:
				pos, ok := ecs.Get[Position](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, pos.z, 3.0)
				hp, ok := ecs.Get[Health](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, hp.hp, 16)
			case npc1:
				pos, ok := ecs.Get[Position](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, pos.z, 1.0)
				hp, ok := ecs.Get[Health](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, hp.hp, 14)
			}
		}
		es = ecs.Query2[Position, Velocity](&world)
		testutil.AssertEqual(t, len(es), 0)
	})

	t.Run("Query3", func(t *testing.T) {
		world := ecs.NewWorld()
		ecs.Initialize[Position](&world)
		ecs.Initialize[Velocity](&world)
		ecs.Initialize[Health](&world)
		player := world.NewEntity()
		npc1 := world.NewEntity()
		npc2 := world.NewEntity()
		testutil.AssertEqual(t, ecs.Add(&world, player, Position{1.0, 2.0, 3.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Velocity{2.0, 2.0, 2.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, player, Health{16}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc1, Position{1.0, 1.0, 1.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc1, Velocity{1.0, 1.0, 1.0}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc1, Health{14}), true)
		testutil.AssertEqual(t, ecs.Add(&world, npc2, Position{2.0, 2.0, 2.0}), true)
		es := ecs.Query3[Position, Velocity, Health](&world)
		testutil.AssertEqual(t, len(es), 2)
		for _, e := range es {
			switch e {
			case player:
				pos, ok := ecs.Get[Position](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, pos.z, 3.0)
				vel, ok := ecs.Get[Velocity](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, vel.x, 2.0)
				hp, ok := ecs.Get[Health](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, hp.hp, 16)
			case npc1:
				pos, ok := ecs.Get[Position](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, pos.z, 1.0)
				vel, ok := ecs.Get[Velocity](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, vel.x, 1.0)
				hp, ok := ecs.Get[Health](&world, e)
				testutil.AssertEqual(t, ok, true)
				testutil.AssertEqual(t, hp.hp, 14)
			}
		}
	})
}

func BenchmarkEcs(b *testing.B) {
	world := ecs.NewWorld()
	ecs.Initialize[Position](&world)
	ecs.Initialize[Velocity](&world)
	ecs.Initialize[Health](&world)
	entities := make([]ecs.Entity, ecs.MAX_ENTITIES)

	for i := range entities {
		entities[i] = world.NewEntity()
		ecs.Add(&world, entities[i], Position{float32(i + 1), 0.0, 0.0})
		if i%64 == 0 {
			ecs.Add(&world, entities[i], Velocity{-float32(i + 1), 0.0, 0.0})
		}
		if i%2 == 0 {
			ecs.Add(&world, entities[i], Health{i + 1})
		}
	}

	b.Logf("Entity Count: %d", world.EntityCount())

	b.Run("Query2", func(b *testing.B) {
		es := ecs.Query2[Position, Health](&world)
		var pcount float32 = 0.0
		hcount := 0
		for _, e := range es {
			pos, _ := ecs.Get[Position](&world, e)
			hp, _ := ecs.Get[Health](&world, e)
			pcount += pos.x
			hcount += hp.hp
		}
	})
}
