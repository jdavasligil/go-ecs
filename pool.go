package ecs

// Pool is a sparse set used for each registered Component type which maps
// entities to their components.
//
// Time Complexity:
//
//	Add    - O(1)
//	Remove - O(1)
//	Query  - O(1)
type Pool[T any] struct {
	// entityIndices is a sparse array that holds the indices into EntityList.
	// The array is indexed by the entity id itself. A value of -1 means empty.
	// TODO: Pagination for memory conservation.
	entityIndices []int

	// entityList is a packed array that contains the entities. The index
	// corresponds to the value from entityIndices.
	entityList []Entity

	// componentList is a packed array that contains component data. The array
	// is aligned with entityList (i.e., entityList[i] corresponds to data in
	// componentList[i]).
	componentList []T
}

// NewPool constructs a component pool for a particular component type.
func NewPool[T any]() Pool[T] {
	p := Pool[T]{
		entityIndices: make([]int, MAX_ENTITIES+1),
		entityList:    make([]Entity, 0, 256),
		componentList: make([]T, 0, 256),
	}
    for i := 0; i < len(p.entityIndices); i++ {
		p.entityIndices[i] = -1
	}
	return p
}

func (p *Pool[T]) IsRegistered(e Entity) bool {
    return p.entityIndices[e.ID()] >= 0
}

// Add registers component of type T to the entity.
func (p *Pool[T]) Add(e Entity, c T) {
    if p.IsRegistered(e) { return }
    p.entityIndices[e.ID()] = len(p.entityList)
    p.entityList = append(p.entityList, e)
    p.componentList = append(p.componentList, c)
}

// Remove unregisters the entity from the component pool.
func (p *Pool[T]) Remove(e Entity) {
    if !p.IsRegistered(e) { return }
    // Get index of the entity to be removed.
    idx := p.entityIndices[e.ID()]
    // Swap the last entity/component with the one marked for removal.
    p.entityList[idx] = p.entityList[len(p.entityList) - 1]
    p.componentList[idx] = p.componentList[len(p.componentList) - 1]
    // Update the new index location for the swapped data.
    p.entityIndices[p.entityList[idx]] = idx
    // Unregister the removed entity.
    p.entityIndices[e.ID()] = -1
    // Delete the last entity/component.
    p.entityList = p.entityList[:len(p.entityList)-1]
    p.componentList = p.componentList[:len(p.componentList)-1]
}

// Retrieves the component data associated with a specific entity.
func (p *Pool[T]) GetComponent(e Entity) (T, bool) {
    if !p.IsRegistered(e) { 
        var noop T
        return noop, false
    }
    return p.componentList[p.entityIndices[e.ID()]], true
}

// Retrieves the slice of all component data independent of the entities.
func (p *Pool[T]) Components() []T {
    return p.componentList
}

// Retrieves the list of all entities with this component registered.
func (p *Pool[T]) Entities() []Entity {
    return p.entityList
}
