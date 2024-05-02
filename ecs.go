package ecs

// https://gist.github.com/dakom/82551fff5d2b843cbe1601bbaff2acbf

type ComponentID uint8

const (
	MAX_ENTITIES   uint32      = 4096 //16777216
	MAX_COMPONENTS ComponentID = 255
)

type World struct {
	entities   EntityManager
	components []any
}

// TODO: How to Query data (entity/components) using Generics?

func Register[T any](e Entity, w *World) {
	w.components = append(w.components, NewComponentStore[T]())
}
