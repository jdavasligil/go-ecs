package ecs_test

import (
	"runtime"
	"testing"

	"github.com/jdavasligil/go-ecs"
	"github.com/jdavasligil/go-ecs/pkg/testutil"
)

func TestQuery(t *testing.T) {
	world := ecs.NewWorld(ecs.WorldOptions{
		EntityLimit:    1024,
		RecycleLimit:   1024,
		ComponentLimit: 255,
	})
	ecs.Initialize[Position](&world)
	ecs.Initialize[Velocity](&world)
	ecs.Initialize[Health](&world)
	ecs.Initialize[CombatTag](&world)
	ecs.Initialize[DeadTag](&world)

	player := world.NewEntity()
	npc1 := world.NewEntity()
	npc2 := world.NewEntity()
	npc3 := world.NewEntity()
	wall := world.NewEntity()

	entIDs := []uint32{
		player.ID(),
		npc1.ID(),
		npc2.ID(),
		npc3.ID(),
		wall.ID(),
	}

	ecs.Add(&world, player, Position{0.0, 0.0, -1.0})
	ecs.Add(&world, player, Velocity{0.0, 0.0, -1.0})
	ecs.Add(&world, player, Health{1})
	ecs.Add(&world, player, CombatTag{})

	ecs.Add(&world, npc1, Position{0.0, 0.0, 1.0})
	ecs.Add(&world, npc1, Health{0})
	ecs.Add(&world, npc1, DeadTag{})

	ecs.Add(&world, npc2, Position{0.0, 0.0, 2.0})
	ecs.Add(&world, npc2, Velocity{0.0, 0.0, 2.0})
	ecs.Add(&world, npc2, Health{3})
	ecs.Add(&world, npc2, CombatTag{})

	ecs.Add(&world, npc3, Position{0.0, 0.0, 3.0})
	ecs.Add(&world, npc3, Velocity{0.0, 0.0, 3.0})
	ecs.Add(&world, npc3, Health{4})

	ecs.Add(&world, wall, Position{0.0, 0.0, 4.0})
	ecs.Add(&world, wall, Health{5})

	t.Run("Get", func(t *testing.T) {
		health, ok := ecs.Get[Health](&world, player)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, health.hp, 1)

		_, dead := ecs.Get[DeadTag](&world, player)
		testutil.AssertEqual(t, dead, false)

		_, cmb := ecs.Get[CombatTag](&world, player)
		testutil.AssertEqual(t, cmb, true)

		pos, ok := ecs.Get[Position](&world, npc1)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, pos.z, 1)
	})

	t.Run("GetMut", func(t *testing.T) {
		healthRef, ok := ecs.GetMut[Health](&world, player)
		testutil.AssertEqual(t, ok, true)
		healthRef.hp = 2

		health, ok := ecs.Get[Health](&world, player)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, health.hp, 2)
		healthRef.hp = 1
	})

	t.Run("Query", func(t *testing.T) {
		es, ps := ecs.Query[Position](&world)
		testutil.AssertEqual(t, len(es), 5)
		testutil.AssertEqual(t, ps[0].x, 0.0)

		seen := make([]bool, world.EntityCount()+1)
		for _, e := range es {
			seen[e.ID()] = true
		}
		for _, id := range entIDs {
			testutil.AssertEqual(t, seen[id], true)
		}
	})

	t.Run("Query2", func(t *testing.T) {
		es := ecs.Query2[Position, CombatTag](&world)
		testutil.AssertEqual(t, len(es), 2)

		seen := make([]bool, world.EntityCount()+1)
		for _, e := range es {
			seen[e.ID()] = true
		}
		testutil.AssertEqual(t, seen[player.ID()], true)
		testutil.AssertEqual(t, seen[npc1.ID()], false)
		testutil.AssertEqual(t, seen[npc2.ID()], true)
		testutil.AssertEqual(t, seen[npc3.ID()], false)
		testutil.AssertEqual(t, seen[wall.ID()], false)

		es = ecs.Query2[CombatTag, DeadTag](&world)
		testutil.AssertEqual(t, len(es), 0)
	})

	t.Run("Query3", func(t *testing.T) {
		es := ecs.Query3[Position, Health, DeadTag](&world)
		testutil.AssertEqual(t, len(es), 1)
		testutil.AssertEqual(t, es[0], npc1)

		es = ecs.Query3[CombatTag, DeadTag, Velocity](&world)
		testutil.AssertEqual(t, len(es), 0)
	})

	t.Run("Query4", func(t *testing.T) {
		es := ecs.Query4[Position, Health, CombatTag, Velocity](&world)
		testutil.AssertEqual(t, len(es), 2)

		seen := make([]bool, world.EntityCount()+1)
		for _, e := range es {
			seen[e.ID()] = true
		}
		testutil.AssertEqual(t, seen[player.ID()], true)
		testutil.AssertEqual(t, seen[npc1.ID()], false)
		testutil.AssertEqual(t, seen[npc2.ID()], true)
		testutil.AssertEqual(t, seen[npc3.ID()], false)
		testutil.AssertEqual(t, seen[wall.ID()], false)

		es = ecs.Query4[CombatTag, DeadTag, Velocity, Health](&world)
		testutil.AssertEqual(t, len(es), 0)
	})

	t.Run("Query5", func(t *testing.T) {
		ecs.Add(&world, player, DeadTag{})
		es := ecs.Query5[Position, Health, CombatTag, Velocity, DeadTag](&world)
		testutil.AssertEqual(t, len(es), 1)
		testutil.AssertEqual(t, es[0], player)

		ecs.Remove[DeadTag](&world, player)

		es = ecs.Query5[CombatTag, DeadTag, Velocity, Health, Position](&world)
		testutil.AssertEqual(t, len(es), 0)
	})
}

var loc int

func BenchmarkQuery(b *testing.B) {
	world := ecs.NewWorld(ecs.WorldOptions{
		EntityLimit:    ecs.MAX_ENTITIES,
		RecycleLimit:   1024,
		ComponentLimit: 255,
	})
	ecs.Initialize[Position](&world)
	ecs.Initialize[Velocity](&world)
	ecs.Initialize[Health](&world)
	ecs.Initialize[CombatTag](&world)
	ecs.Initialize[DeadTag](&world)

	entities := make([]ecs.Entity, world.EntityLimit())

	for i := range entities {
		entities[i] = world.NewEntity()
		ecs.Add(&world, entities[i], Position{float32(i + 1), 0.0, 0.0})
		if i%2 == 0 {
			ecs.Add(&world, entities[i], Health{i + 1})
		}
		if i%256 == 0 {
			ecs.Add(&world, entities[i], Velocity{-float32(i + 1), 0.0, 0.0})
		}
		if i%512 == 0 {
			ecs.Add(&world, entities[i], CombatTag{})
		}
		if i%1024 == 0 {
			ecs.Add(&world, entities[i], DeadTag{})
		}
	}

	b.Logf("Entity Count: %d", world.EntityCount())
	b.Log("Memory Usage:")
	b.Logf("    World     - %d", world.MemUsage())
	b.Logf("    Position  - %d", ecs.MemUsage[Position](&world))
	b.Logf("    Velocity  - %d", ecs.MemUsage[Velocity](&world))
	b.Logf("    Health    - %d", ecs.MemUsage[Health](&world))
	b.Logf("    CombatTag - %d", ecs.MemUsage[CombatTag](&world))
	b.Logf("    DeadTag   - %d", ecs.MemUsage[DeadTag](&world))
	b.Logf("Total: %d",
		world.MemUsage()+
			ecs.MemUsage[Position](&world)+
			ecs.MemUsage[Velocity](&world)+
			ecs.MemUsage[Health](&world)+
			ecs.MemUsage[CombatTag](&world)+
			ecs.MemUsage[DeadTag](&world))

	b.ResetTimer()
	b.Run("Query2", func(b *testing.B) {
		b.ReportAllocs()
		b.StartTimer()
		es := ecs.Query2[Position, Health](&world)
		for _, e := range es {
			hp, _ := ecs.Get[Health](&world, e)
			loc += hp.hp
		}
		b.StopTimer()
	})
	b.Run("Query5", func(b *testing.B) {
		b.StartTimer()
		es := ecs.Query4Exclude1[Position, Velocity, Health, CombatTag, DeadTag](&world)
		for _, e := range es {
			hp, _ := ecs.Get[Health](&world, e)
			loc += hp.hp
		}
		b.StopTimer()
	})
	b.Run("QueryExclude", func(b *testing.B) {
		b.StartTimer()
		es := ecs.QueryExclude[Position, Health](&world)
		b.StopTimer()
		for _, e := range es {
			hp, ok := ecs.Get[Health](&world, e)
			if ok {
				loc += hp.hp
			}
		}
	})
	runtime.KeepAlive(loc)
}
