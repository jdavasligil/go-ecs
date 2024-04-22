package ecs

// https://gist.github.com/dakom/82551fff5d2b843cbe1601bbaff2acbf

type ComponentID uint8

const (
	MAX_ENTITIES   Entity      = 16777216
	MAX_COMPONENTS ComponentID = 255
)

type World struct {
}
