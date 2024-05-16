package ecs

type ComponentID uint8

// Component represents any pure data type and is given an ID.
type Component interface {
	ID() ComponentID
}

// Initialize initializes a component which ensures that a store is created.
//
// Initialize must be called before entities are added.
func Initialize[T Component](w *World) bool {
	var noop T
	if w.components[noop.ID()] != nil || w.ComponentCount == cap(w.components) {
		return false
	}
	w.components[noop.ID()] = newComponentStore[T]()
	w.ComponentCount++
	return true
}

// Add adds a component to an entity if that component was initialized.
func Add[T Component](w *World, e Entity, c T) bool {
	store, ok := w.components[c.ID()].(*componentStore[T])
	if !ok {
		return false
	}
	return store.Add(e, c)
}

// Remove removes a component from an entity.
//
// Remove opts out of cleaning unused page memory for peformance.
// To clean paged memory use RemoveAndClean.
//
// Note that removal does not preserve order in the packed arrays.
//
// Time Complexity: O(1)
func Remove[T Component](w *World, e Entity) bool {
	var noop T
	store, ok := w.components[noop.ID()].(*componentStore[T])
	if !ok {
		return ok
	}
	return store.Remove(e)
}

// RemoveAndClean removes a component from an entity and sweeps the page.
//
// Additionally, it performs a reallocation of dense arrays which may cause
// GC pauses if overused.
//
// Sweeped pages are only cleaned up if completely empty. Hence, the majority
// of the time it does nothing. Sweeping is ~32x slower. For performance
// critical tasks use Remove.
//
// Time Complexity: O(N) where N is the page size.
func RemoveAndClean[T Component](w *World, e Entity) bool {
	var noop T
	store, ok := w.components[noop.ID()].(*componentStore[T])
	if !ok {
		return ok
	}
	return store.RemoveAndClean(e)
}

// Sweep iterates through the component store freeing memory of empty pages.
//
// It is not a bad idea to call this periodically if removals are frequent
// and space is a premium.
//
// Time Complexity: O(MN) where M is the page count and N is the page size.
func Sweep[T Component](w *World) {
	var noop T
	store, ok := w.components[noop.ID()].(*componentStore[T])
	if !ok {
		return
	}
	store.entityIndices.Sweep()
}

// MemUsage reports the memory being used by the component store in bytes.
func MemUsage[T Component](w *World) uintptr {
	var noop T
	store, ok := w.components[noop.ID()].(*componentStore[T])
	if !ok {
		return 0
	}
	return store.MemUsage()
}
