package main

import (
	"fmt"

	"github.com/jdavasligil/go-ecs"
)

// Each component needs a unique ID. Creating an enum for all of your
// components is an easy way to accomplish that.
const (
	PositionID ecs.ComponentID = iota
	VelocityID
	TagID
)

// Components are just data with an associated ID.
type Position struct {
	x, y float32
}
type Velocity struct {
	x, y float32
}

// Tags are components with no data. Useful for filtering queries.
type Tag struct{}

// Every component must have an ID Method that returns a ComponentID.
// This is used by the ECS to index the internal component store.
func (c Position) ID() ecs.ComponentID { return PositionID }
func (c Velocity) ID() ecs.ComponentID { return VelocityID }
func (c Tag) ID() ecs.ComponentID      { return TagID }

func main() {
	// The world tracks entities and their components. It can be passed around
	// by reference for any system to have access.
	world := ecs.NewWorld(ecs.WorldOptions{
		EntityLimit:    10_000,             // Limit on living entities.
		RecycleLimit:   1_000,              // Limit on recycled entities.
		ComponentLimit: ecs.MAX_COMPONENTS, // Limit on components registered.
	})

	// Component stores must be initialized for each component.
	ecs.Initialize[Position](&world)
	ecs.Initialize[Velocity](&world)
	ecs.Initialize[Tag](&world)

	entity1 := world.NewEntity()
	entity2 := world.NewEntity()

	ecs.Add(&world, entity1, Position{})
	ecs.Add(&world, entity1, Velocity{})
	ecs.Add(&world, entity1, Tag{})

	ecs.Add(&world, entity2, Position{})
	ecs.Add(&world, entity2, Velocity{})

	// Querying a single component gives you mutable access to the list of
	// components and entities with those components in alignment.
	es, ps := ecs.Query[Position](&world)
	_, vs := ecs.Query[Velocity](&world)

	// Alignment between multiple queries is not guaranteed by the ECS.
	// However, it is possible if additions and removals are synchronized.
	for i, e := range es {
		ps[i].y += 1.0
		vs[i].x += -0.5
		fmt.Printf("Entity ID: %d\n", e.ID())
		fmt.Printf("    Position - %v\n", ps[i])
		fmt.Printf("    Velocity - %v\n", vs[i])
	}

	// Queries for the intersection of multiple components can be performed.
	// This only returns the list of entities which share N components.
	es = ecs.Query3[Position, Velocity, Tag](&world)

	for _, e := range es {
		// Accessing the components can be done using Get or GetMut.
		// Use Get if you want a copy. GetMut returns a pointer to the data.
		p, _ := ecs.GetMut[Position](&world, e)
		v, _ := ecs.GetMut[Velocity](&world, e)
		p.x += 1.0
		v.y -= 0.5
		fmt.Printf("Entity ID: %d\n", e.ID())
		fmt.Printf("    Position - %v\n", p)
		fmt.Printf("    Velocity - %v\n", v)
	}

	// Removal must be done manually for each component.
	ecs.Remove[Position](&world, entity2)
	ecs.Remove[Velocity](&world, entity2)

	// Remove and clean performs a reallocation preventing a memory leak.
	// This only needs to be performed once per component store at the end.
	ecs.RemoveAndClean[Position](&world, entity1)
	ecs.RemoveAndClean[Velocity](&world, entity1)
	ecs.RemoveAndClean[Tag](&world, entity1)

	// Once the components are removed, its safe to destroy the entities.
	world.DestroyEntity(entity1)
	world.DestroyEntity(entity2)

	// Note that entity IDs are recycled. We now have dangling references.
	entity1v2 := world.NewEntity()
	fmt.Printf("Old ID: %d, Version: %d\n", entity1.ID(), entity1.Version())
	fmt.Printf("New ID: %d, Version: %d\n", entity1v2.ID(), entity1v2.Version())

	// The version increments when an entity is destroyed. This means you can
	// validate your entity reference by comparing versions.
	//
	// This rolls over after 255 generations. Hence, it is still possible to
	// have an incorrect match, though it is unlikely.
	//
	// I would recommend being careful when holding on to entity references
	// after deletion. An event system could be used to alert any subscribers
	// that its reference is no longer valid.
}
