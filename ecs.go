package ecs

// https://gist.github.com/dakom/82551fff5d2b843cbe1601bbaff2acbf

const (
	MAX_ENTITIES   uint32      = 8847360 // 16777216 Theoretical Max
	MAX_COMPONENTS ComponentID = 255
)

// World contains all entities and their components. ECS provides an API for
// creating systems which can query and act upon the packed component data.
type World struct {
	entities       EntityManager
	components     []any
	ComponentCount int
}

func NewWorld() World {
	return World{
		entities:   NewEntityManager(),
		components: make([]any, MAX_COMPONENTS),
	}
}

func (w *World) NewEntity() Entity {
	return w.entities.CreateEntity()
}

// DestroyEntity only recycles the associated Entity ID. It is up to the caller
// to ensure that the entity is removed from each component store using the
// Remove function. Otherwise, this creates a memory leak.
func (w *World) DestroyEntity(e Entity) bool {
	return w.entities.RecycleEntity(e)
}

func (w *World) EntityCount() int {
	return int(w.entities.size)
}

// Initialize is used to initialize a component which ensures that a store
// is created for that component. This is a performance optimization.
func Initialize[T Component](w *World) bool {
	var noop T
	if w.components[noop.ID()] != nil {
		return false
	}
	w.components[noop.ID()] = NewComponentStore[T]()
	w.ComponentCount++
	return true
}

// Register is used to register component data with an entity. This function
// is used to both initialize and add component data to an entity.
func Register[T Component](w *World, e Entity, c T) bool {
	var store *ComponentStore[T]
	if w.components[c.ID()] == nil {
		store = NewComponentStore[T]()
		w.components[c.ID()] = store
		w.ComponentCount++
	} else {
		var ok bool
		store, ok = w.components[c.ID()].(*ComponentStore[T])
		if !ok {
			return false
		}
	}
	return store.Add(e, c)
}

// Add will add a component to an entity.
func Add[T Component](w *World, e Entity, c T) bool {
	store, ok := w.components[c.ID()].(*ComponentStore[T])
	if !ok {
		return false
	}
	return store.Add(e, c)
}

// Remove will remove a component from an entity and returns the affected
// entity or the null entity upon failure.
func Remove[T Component](w *World, e Entity) Entity {
	var noop T
	store, ok := w.components[noop.ID()].(*ComponentStore[T])
	if !ok {
		return 0
	}
	return store.Remove(e)
}

// Query returns a copy of the data queried for a single entity.
func Query[T Component](w *World, e Entity) (T, bool) {
	var noop T
	store, ok := w.components[noop.ID()].(*ComponentStore[T])
	if !ok {
		return noop, ok
	}
	return store.GetComponent(e)
}

// MutQuery returns a mutable reference to the underlying data queried for
// a particular entity. Only a single caller may claim ownership at a time.
func MutQuery[T Component](w *World, e Entity) (*T, bool) {
	var noop T
	store, ok := w.components[noop.ID()].(*ComponentStore[T])
	if !ok {
		return nil, ok
	}
	return store.GetMutComponent(e)
}

// QueryAll returns slices to both the entities and their underlying data. The
// data is mutable, packed, aligned, and so can be iterated together. Only a
// single caller may claim mutable ownership at a time. Slices are possibly nil.
func QueryAll[T Component](w *World) ([]Entity, []T) {
	var noop T
	store, ok := w.components[noop.ID()].(*ComponentStore[T])
	if !ok {
		return nil, nil
	}
	return store.entityList, store.componentList
}

// QueryIntersect2 performs a query finding the intersection of two components.
// It returns the aligned slices of entities which have both components and
// their associated data in order. Slices are possibly nil.
//
// Time Complexity: O(N) where N = min(# Entities with T, # Entities with V)
func QueryIntersect2[T Component, V Component](w *World) ([]Entity, []T, []V) {
	var noopT T
	var noopV V
	storeT, okT := w.components[noopT.ID()].(*ComponentStore[T])
	storeV, okV := w.components[noopV.ID()].(*ComponentStore[V])
	if !(okT && okV) {
		return nil, nil, nil
	}
	if len(storeT.entityList) < len(storeV.entityList) {
		es := make([]Entity, 0)
		ts := make([]T, 0)
		vs := make([]V, 0)
		for idxT, e := range storeT.entityList {
			idxV := storeV.entityIndices.At(int(e.ID()))
			if idxV >= 0 {
				es = append(es, e)
				ts = append(ts, storeT.componentList[idxT])
				vs = append(vs, storeV.componentList[idxV])
			}
		}
		return es, ts, vs
	} else {
		es := make([]Entity, 0)
		ts := make([]T, 0)
		vs := make([]V, 0)
		for idxV, e := range storeV.entityList {
			idxT := storeT.entityIndices.At(int(e.ID()))
			if idxT >= 0 {
				es = append(es, e)
				ts = append(ts, storeT.componentList[idxT])
				vs = append(vs, storeV.componentList[idxV])
			}
		}
		return es, ts, vs
	}
}

// QueryEntities2 performs a query finding the intersection of two components.
// It returns a packed slice of entities which have both components but not
// their associated data. Used with QueryMut to mutate data. Slices can be nil.
//
// Time Complexity: O(N) where N = min(# Entities with T, # Entities with V)
func QueryEntities2[T Component, V Component](w *World) []Entity {
	var noopT T
	var noopV V
	storeT, okT := w.components[noopT.ID()].(*ComponentStore[T])
	storeV, okV := w.components[noopV.ID()].(*ComponentStore[V])
	if !(okT && okV) {
		return nil
	}
	es := make([]Entity, 0)
	if len(storeT.entityList) < len(storeV.entityList) {
		for _, e := range storeT.entityList {
			idxV := storeV.entityIndices.At(int(e.ID()))
			if idxV >= 0 {
				es = append(es, e)
			}
		}
	} else {
		for _, e := range storeV.entityList {
			idxT := storeT.entityIndices.At(int(e.ID()))
			if idxT >= 0 {
				es = append(es, e)
			}
		}
	}
	return es
}

// QueryEntities3 performs a query for the intersection of three components.
// It returns a packed slice of entities which have all components but not
// their associated data. Used with QueryMut to mutate data. Slices can be nil.
//
// Time Complexity: O(N) where N = min(# Entities of a component type)
func QueryEntities3[A Component, B Component, C Component](w *World) []Entity {
	var noopA A
	var noopB B
	var noopC C
	storeA, okA := w.components[noopA.ID()].(*ComponentStore[A])
	storeB, okB := w.components[noopB.ID()].(*ComponentStore[B])
	storeC, okC := w.components[noopC.ID()].(*ComponentStore[C])
	if !(okA && okB && okC) {
		return nil
	}
	es := make([]Entity, 0)
	lenA := len(storeA.entityList)
	lenB := len(storeB.entityList)
	lenC := len(storeC.entityList)
	if (lenA <= lenB && lenB <= lenC) ||
		(lenA <= lenC && lenC <= lenB) {
		for _, e := range storeA.entityList {
			idxB := storeB.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			if idxB >= 0 &&
				idxC >= 0 {
				es = append(es, e)
			}
		}
	} else if (lenB <= lenA && lenA <= lenC) ||
		(lenB <= lenC && lenC <= lenA) {
		for _, e := range storeB.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxC >= 0 {
				es = append(es, e)
			}
		}
	} else {
		for _, e := range storeC.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxB := storeB.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxB >= 0 {
				es = append(es, e)
			}
		}
	}
	return es
}
