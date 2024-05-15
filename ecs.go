// Package ecs is an Entity Component System library. This implementation uses
// a sparse set design optimized for flexibility in quickly adding and removing
// components for highly dynamic games.
//
// This library simply provides a way to manage entities and their components
// through the World struct and then query those components. Serialization,
// scheduling, event handling etc. is out of the scope for this library.
//
// Design is heavily inspired by the research done by dakom on EnTT & Shipyard.
// https://gist.github.com/dakom/82551fff5d2b843cbe1601bbaff2acbf
//
// go-ecs is licensed under the dual MIT/APACHE 2.0 license.

package ecs

import "unsafe"

const (
	MAX_ENTITIES   uint32      = 16777216
	MAX_COMPONENTS ComponentID = 255
)

// World contains all entities and their components.
type World struct {
	entities       entityManager
	components     []any
	ComponentCount int
}

// WorldOptions lists the option parameters required to create a World.
type WorldOptions struct {
	// EntityLimit is the maximum number of living entities permitted.
	// Bounded by MAX_ENTITIES.
	EntityLimit uint32

	// RecycleLimit is the maximum number of dead entities kept for recycling.
	// Bounded by EntityLimit.
	RecycleLimit uint32

	// ComponentLimit is the maximum number of component types.
	// Bounded by MAX_COMPONENTS.
	ComponentLimit ComponentID
}

// NewWorld creates a new world with the given options.
func NewWorld(opts WorldOptions) World {
	return World{
		entities:   newEntityManager(opts.EntityLimit, opts.RecycleLimit),
		components: make([]any, opts.ComponentLimit),
	}
}

// NewEntity creates a new Entity.
//
// The null entity is returned upon failure.
func (w *World) NewEntity() Entity {
	return w.entities.CreateEntity()
}

// DestroyEntity recycles the associated Entity ID.
//
// Warning: It is up to the caller to ensure that the entity is removed from
// each component store using the remove function to prevent a memory leak.
func (w *World) DestroyEntity(e Entity) bool {
	return w.entities.RecycleEntity(e)
}

func (w *World) EntityCount() int {
	return int(w.entities.size)
}

func (w *World) EntityLimit() int {
	return int(w.entities.MaxEntities)
}

func (w *World) RecycleLimit() int {
	return int(w.entities.MaxRecycle)
}

func (w *World) ComponentLimit() int {
	return cap(w.components)
}

// MemUsage for the world does not include the memory taken by the component
// stores. The MemUsage of each component store must be added for a total.
func (w *World) MemUsage() uintptr {
	size := unsafe.Sizeof(*w)
	size += w.entities.MemUsage()
	size += unsafe.Sizeof(w.components)
	size += unsafe.Sizeof(w.ComponentCount)
	return size
}
