package ecs

import "github.com/jdavasligil/go-ecs/pkg/bitset"

type Archetype uint32

type ComponentType uint8

type Signature = bitset.BitsetUint64

const (
	MAX_ENTITIES   Entity        = 1e6
	MAX_COMPONENTS ComponentType = 255
)

type World struct {
}
