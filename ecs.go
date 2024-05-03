package ecs

// https://gist.github.com/dakom/82551fff5d2b843cbe1601bbaff2acbf

const (
	MAX_ENTITIES   uint32      = 4096 //16777216
	MAX_COMPONENTS ComponentID = 255
)

// CONCURRENCY: thread safety on component access?
type World struct {
	entities   EntityManager
	components []any
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

func (w *World) EntityCount() int {
    return int(w.entities.size)
}

func (w *World) ComponentCount() int {
    count := 0
    for _, c := range w.components {
        if c != nil {
            count++
        }
    }
	return count
}

// Register is used to register component data with an entity.
func Register[T Component](w *World, e Entity, c T) bool {
	if w.components[c.ID()] == nil {
		store := NewComponentStore[T]()
		store.Add(e, c)
		w.components[c.ID()] = store
	} else {
		store, ok := w.components[c.ID()].(*ComponentStore[T])
		if !ok {
			return false
		}
		store.Add(e, c)
	}
	return true
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
// single caller may claim mutable ownership at a time.
func QueryAll[T Component](w *World) ([]Entity, []T) {
	var noop T
	store, ok := w.components[noop.ID()].(*ComponentStore[T])
	if !ok {
		return nil, nil
	}
    return store.entityList, store.componentList
}

func Query2[T Component, V Component](w *World) ([]Entity, []T, []V) {
	var noopT T
	var noopV V
	storeT, okT := w.components[noopT.ID()].(*ComponentStore[T])
	storeV, okV := w.components[noopV.ID()].(*ComponentStore[V])
	if !(okT || okV) {
        if okT {
            return storeT.entityList, storeT.componentList, nil
        } else if okV {
            return storeV.entityList, nil, storeV.componentList
        } else {
            return nil, nil, nil
        }
	}
    if len(storeT.entityList) < len(storeV.entityList) {
        es := make([]Entity, 0)
        ts := make([]T, 0)
        vs := make([]V, 0)
        for idxT, e := range storeT.entityList {
            idxV := storeV.entityIndices.At(int(e.ID()))
            if  idxV >= 0 {
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
            if  idxT >= 0 {
                es = append(es, e)
                ts = append(ts, storeT.componentList[idxT])
                vs = append(vs, storeV.componentList[idxV])
            }
        }
        return es, ts, vs
    }
}
