package ecs

import (
	"unsafe"

	"github.com/jdavasligil/go-ecs/pkg/pagearray"
)

// componentStore is a sparse set used for each registered Component type which maps
// entities to their components.
//
// Time Complexity:
//
//	Add    - O(1)
//	Remove - O(1)
//	Query  - O(1)
type componentStore[T Component] struct {
	// entityIndices is a sparse array that holds the indices into EntityList.
	// The array is indexed by the entity id itself. A value of -1 means empty.
	entityIndices pagearray.PageArray

	// entityList is a packed array that contains the entities. The index
	// corresponds to the value from entityIndices.
	entityList []Entity

	// componentList is a packed array that contains component data. The array
	// is aligned with entityList (i.e., entityList[i] corresponds to data in
	// componentList[i]).
	componentList []T
}

// NewcomponentStore constructs a component store for a particular component type.
func NewComponentStore[T Component]() *componentStore[T] {
	p := &componentStore[T]{
		entityIndices: pagearray.NewPageArray(10), // Page size 1024
		entityList:    make([]Entity, 0),
		componentList: make([]T, 0),
	}
	return p
}

func (p *componentStore[T]) IsRegistered(e Entity) bool {
	return p.entityIndices.At(int(e.ID())) >= 0
}

// Add registers component of type T to the entity. Returns true if successful.
func (p *componentStore[T]) Add(e Entity, c T) bool {
	if p.IsRegistered(e) {
		return false
	}
	p.entityIndices.Set(int(e.ID()), len(p.entityList))
	p.entityList = append(p.entityList, e)
	p.componentList = append(p.componentList, c)
	return true
}

// Remove unregisters the entity from the component store and returns the
// removed entity. Page memory is not cleaned.
func (p *componentStore[T]) RemoveAndClean(e Entity) Entity {
	if !p.IsRegistered(e) {
		return 0
	}
	// Get index of the entity to be removed.
	idx := p.entityIndices.At(int(e.ID()))
	entity := p.entityList[idx]
	// Swap the last entity/component with the one marked for removal.
	p.entityList[idx] = p.entityList[len(p.entityList)-1]
	p.componentList[idx] = p.componentList[len(p.componentList)-1]
	// Update the new index location for the swapped data.
	p.entityIndices.Set(int(p.entityList[idx].ID()), idx)
	// Unregister the removed entity.
	p.entityIndices.SweepAndClear(int(e.ID()))
	// Delete the last entity/component.
	p.entityList = p.entityList[:len(p.entityList)-1]
	p.componentList = p.componentList[:len(p.componentList)-1]

	return entity
}

// Remove unregisters the entity from the component store and returns the
// removed entity. Page memory is not cleaned.
func (p *componentStore[T]) Remove(e Entity) Entity {
	if !p.IsRegistered(e) {
		return 0
	}
	// Get index of the entity to be removed.
	idx := p.entityIndices.At(int(e.ID()))
	entity := p.entityList[idx]
	// Swap the last entity/component with the one marked for removal.
	p.entityList[idx] = p.entityList[len(p.entityList)-1]
	p.componentList[idx] = p.componentList[len(p.componentList)-1]
	// Update the new index location for the swapped data.
	p.entityIndices.Set(int(p.entityList[idx].ID()), idx)
	// Unregister the removed entity.
	p.entityIndices.Clear(int(e.ID()))
	// Delete the last entity/component.
	p.entityList = p.entityList[:len(p.entityList)-1]
	p.componentList = p.componentList[:len(p.componentList)-1]

	return entity
}

// Retrieves the component data associated with a specific entity.
func (p *componentStore[T]) GetComponent(e Entity) (T, bool) {
	if !p.IsRegistered(e) {
		var noop T
		return noop, false
	}
	return p.componentList[p.entityIndices.At(int(e.ID()))], true
}

// Retrieves a mutable reference to the  component data associated with a
// specific entity. Only a single caller may claim ownership at a time.
func (p *componentStore[T]) GetMutComponent(e Entity) (*T, bool) {
	if !p.IsRegistered(e) {
		return nil, false
	}
	return &p.componentList[p.entityIndices.At(int(e.ID()))], true
}

func (p *componentStore[T]) Entities() []Entity {
	return p.entityList
}

func (p *componentStore[T]) Components() []T {
	return p.componentList
}

func (p *componentStore[T]) Size() int {
	return len(p.entityList)
}

// Reset performs a hard reset by throwing away all allocated memory for
// garbage collection. May negatively affect garbage collection performance.
func (p *componentStore[T]) Reset() {
	p.entityIndices.Reset()
	p.entityList = make([]Entity, 0, 256)
	p.componentList = make([]T, 0, 256)
}

// MemUsage returns an estimate for the current memory being used in bytes.
func (p *componentStore[T]) MemUsage() uintptr {
	var entityType Entity
	var componentType T
	size := unsafe.Sizeof(*p)
	size += unsafe.Sizeof(p.entityList)
	size += unsafe.Sizeof(p.componentList)
	size += unsafe.Sizeof(entityType) * uintptr(cap(p.entityList))
	size += unsafe.Sizeof(componentType) * uintptr(cap(p.componentList))
	size += p.entityIndices.MemUsage()
	return size
}
