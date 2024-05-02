package ecs

import (
	"sync"

	"github.com/jdavasligil/go-ecs/pkg/pagearray"
)

// Pool is a sparse set used for each registered Component type which maps
// entities to their components.
//
// Time Complexity:
//
//	Add    - O(1)
//	Remove - O(1)
//	Query  - O(1)
//
// Thread Safety: PageArray is thread safe, but Pool is not. To protect pool
// from concurrent access, Pool may be locked and unlocked directly. E.g.,
//
// pool := NewPool[Component]()
// pool.Lock()
// pool.Add(entity, component)
// pool.Unlock()
type Pool[T any] struct {
	sync.Mutex

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

// NewPool constructs a component pool for a particular component type.
func NewPool[T any]() *Pool[T] {
	p := &Pool[T]{
		entityIndices: pagearray.NewPageArray(10), // Page size 1024
		entityList:    make([]Entity, 0, 256),
		componentList: make([]T, 0, 256),
	}
	return p
}

func (p *Pool[T]) IsRegistered(e Entity) bool {
	return p.entityIndices.At(int(e.ID())) >= 0
}

// Add registers component of type T to the entity. Returns true if successful.
func (p *Pool[T]) Add(e Entity, c T) bool {
	if p.IsRegistered(e) {
		return false
	}
	p.entityIndices.Set(int(e.ID()), len(p.entityList))
	p.entityList = append(p.entityList, e)
	p.componentList = append(p.componentList, c)
	return true
}

// Remove unregisters the entity from the component pool and returns the
// removed entity.
func (p *Pool[T]) Remove(e Entity) Entity {
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

// Retrieves the component data associated with a specific entity.
func (p *Pool[T]) GetComponent(e Entity) (T, bool) {
	if !p.IsRegistered(e) {
		var noop T
		return noop, false
	}
	return p.componentList[p.entityIndices.At(int(e.ID()))], true
}

// Retrieves the slice of all component data independent of the entities.
func (p *Pool[T]) Components() []T {
	return p.componentList
}

// Retrieves the list of all entities with this component registered.
func (p *Pool[T]) Entities() []Entity {
	return p.entityList
}

func (p *Pool[T]) Size() int {
	return len(p.entityList)
}

// Reset performs a hard reset by throwing away all allocated memory for
// garbage collection. May negatively affect garbage collection performance.
func (p *Pool[T]) Reset() {
	p.entityIndices.Reset()
	p.entityList = make([]Entity, 0, 256)
	p.componentList = make([]T, 0, 256)
}
