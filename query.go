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

// Query2 performs a query finding the intersection of two components.
//
// It returns a packed slice of entities which have both components but not
// their associated data. Paired with GetMut to mutate data.
//
// Time Complexity: O(N) where N = min(# Entities with T, # Entities with V)
func Query2[T Component, V Component](w *World) []Entity {
	var noopT T
	var noopV V
	storeT, okT := w.components[noopT.ID()].(*componentStore[T])
	storeV, okV := w.components[noopV.ID()].(*componentStore[V])
	es := make([]Entity, 0)
	if !(okT && okV) {
		return es
	}
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

// Query3 performs a query for the intersection of three components.
//
// It returns a packed slice of entities which have all components but not
// their associated data. Paired with GetMut to mutate data.
//
// Time Complexity: O(N) where N = min(# Entities of a component type)
func Query3[A Component, B Component, C Component](w *World) []Entity {
	var noopA A
	var noopB B
	var noopC C
	storeA, okA := w.components[noopA.ID()].(*componentStore[A])
	storeB, okB := w.components[noopB.ID()].(*componentStore[B])
	storeC, okC := w.components[noopC.ID()].(*componentStore[C])
	es := make([]Entity, 0)
	if !(okA && okB && okC) {
		return es
	}
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

// Query4 performs a query for the intersection of four components.
//
// It returns a packed slice of entities which have all components but not
// their associated data. Paired with GetMut to mutate data.
//
// Time Complexity: O(N) where N = min(# Entities of a component type)
func Query4[
	A Component,
	B Component,
	C Component,
	D Component,
](w *World) []Entity {
	var noopA A
	var noopB B
	var noopC C
	var noopD D
	storeA, okA := w.components[noopA.ID()].(*componentStore[A])
	storeB, okB := w.components[noopB.ID()].(*componentStore[B])
	storeC, okC := w.components[noopC.ID()].(*componentStore[C])
	storeD, okD := w.components[noopD.ID()].(*componentStore[D])
	es := make([]Entity, 0)
	if !(okA && okB && okC && okD) {
		return es
	}
	lenA := len(storeA.entityList)
	lenB := len(storeB.entityList)
	lenC := len(storeC.entityList)
	lenD := len(storeD.entityList)
	minLen := min(lenA, lenB, lenC, lenD)
	switch minLen {
	case lenA:
		for _, e := range storeA.entityList {
			idxB := storeB.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			idxD := storeD.entityIndices.At(int(e.ID()))
			if idxB >= 0 &&
				idxC >= 0 &&
				idxD >= 0 {
				es = append(es, e)
			}
		}
	case lenB:
		for _, e := range storeB.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			idxD := storeD.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxC >= 0 &&
				idxD >= 0 {
				es = append(es, e)
			}
		}
	case lenC:
		for _, e := range storeC.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxB := storeB.entityIndices.At(int(e.ID()))
			idxD := storeD.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxB >= 0 &&
				idxD >= 0 {
				es = append(es, e)
			}
		}
	case lenD:
		for _, e := range storeD.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxB := storeB.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxB >= 0 &&
				idxC >= 0 {
				es = append(es, e)
			}
		}
	}
	return es
}

// Query5 performs a query for the intersection of five components.
//
// It returns a packed slice of entities which have all components but not
// their associated data. Paired with GetMut to mutate data.
//
// Time Complexity: O(N) where N = min(# Entities of a component type)
func Query5[
	A Component,
	B Component,
	C Component,
	D Component,
	E Component,
](w *World) []Entity {
	var noopA A
	var noopB B
	var noopC C
	var noopD D
	var noopE E
	storeA, okA := w.components[noopA.ID()].(*componentStore[A])
	storeB, okB := w.components[noopB.ID()].(*componentStore[B])
	storeC, okC := w.components[noopC.ID()].(*componentStore[C])
	storeD, okD := w.components[noopD.ID()].(*componentStore[D])
	storeE, okE := w.components[noopE.ID()].(*componentStore[E])
	es := make([]Entity, 0)
	if !(okA && okB && okC && okD && okE) {
		return es
	}
	lenA := len(storeA.entityList)
	lenB := len(storeB.entityList)
	lenC := len(storeC.entityList)
	lenD := len(storeD.entityList)
	lenE := len(storeE.entityList)
	minLen := min(lenA, lenB, lenC, lenD, lenE)
	switch minLen {
	case lenA:
		for _, e := range storeA.entityList {
			idxB := storeB.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			idxD := storeD.entityIndices.At(int(e.ID()))
			idxE := storeE.entityIndices.At(int(e.ID()))
			if idxB >= 0 &&
				idxC >= 0 &&
				idxD >= 0 &&
				idxE >= 0 {
				es = append(es, e)
			}
		}
	case lenB:
		for _, e := range storeB.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			idxD := storeD.entityIndices.At(int(e.ID()))
			idxE := storeE.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxC >= 0 &&
				idxD >= 0 &&
				idxE >= 0 {
				es = append(es, e)
			}
		}
	case lenC:
		for _, e := range storeC.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxB := storeB.entityIndices.At(int(e.ID()))
			idxD := storeD.entityIndices.At(int(e.ID()))
			idxE := storeE.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxB >= 0 &&
				idxD >= 0 &&
				idxE >= 0 {
				es = append(es, e)
			}
		}
	case lenD:
		for _, e := range storeD.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxB := storeB.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			idxE := storeE.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxB >= 0 &&
				idxC >= 0 &&
				idxE >= 0 {
				es = append(es, e)
			}
		}
	case lenE:
		for _, e := range storeE.entityList {
			idxA := storeA.entityIndices.At(int(e.ID()))
			idxB := storeB.entityIndices.At(int(e.ID()))
			idxC := storeC.entityIndices.At(int(e.ID()))
			idxD := storeD.entityIndices.At(int(e.ID()))
			if idxA >= 0 &&
				idxB >= 0 &&
				idxC >= 0 &&
				idxD >= 0 {
				es = append(es, e)
			}
		}
	}
	return es
}
