package ecs

// Get returns a copy of the component for a single entity.
func Get[T Component](w *World, e Entity) (T, bool) {
	var noop T
	store, ok := w.components[noop.ID()].(*componentStore[T])
	if !ok {
		return noop, ok
	}
	return store.GetComponent(e)
}

// GetMut returns a mutable reference to the underlying component for
// a single entity. Only a single caller may claim ownership at a time.
//
// Reference is possibly nil.
func GetMut[T Component](w *World, e Entity) (*T, bool) {
	var noop T
	store, ok := w.components[noop.ID()].(*componentStore[T])
	if !ok {
		return nil, ok
	}
	return store.GetMutComponent(e)
}

// Query returns slices to both the entities and their underlying data.
//
// The data is mutable, packed, aligned, and so can be iterated together. Only a
// single caller may claim mutable ownership at a time.
//
// Slices are possibly nil.
func Query[T Component](w *World) ([]Entity, []T) {
	var noop T
	store, ok := w.components[noop.ID()].(*componentStore[T])
	if !ok {
		return nil, nil
	}
	return store.entityList, store.componentList
}

// QueryExclude performs a query finding the set difference of component T\V.
//
// It returns a packed slice of entities which have component T but not
// component V. Paired with GetMut to mutate data.
//
// Time Complexity: O(N) where N = # Entities with T
func QueryExclude[T Component, V Component](w *World) []Entity {
	var noopT T
	var noopV V
	storeT, okT := w.components[noopT.ID()].(*componentStore[T])
	storeV, okV := w.components[noopV.ID()].(*componentStore[V])
	es := make([]Entity, 0)
	if !(okT && okV) {
		return es
	}
	for _, e := range storeT.entityList {
		if storeV.entityIndices.At(int(e.ID())) < 0 {
			es = append(es, e)
		}
	}
	return es
}

//go:generate go run internal/gen/query_gen.go -N=6
