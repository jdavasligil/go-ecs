package ecs

// [ID                    ][Version]
// ########################VVVVVVVV
// 00000000000000000000000011111111

// Entity is a unique ID which corresponds to exactly one game object.
// It is used to reference a collection of Components (data).
type Entity uint32

// NewEntity requires that the provided id < 16777215.
func newEntity(id uint32) Entity {
	return Entity(id << 8)
}

// Version returns the current generation of the entity. Used for comparison
// to verify that an existing reference is invalid (entity was deleted).
func (e *Entity) Version() uint8 {
	return uint8(*e)
}

// ID returns the actual unique ID of the Entity. This shall be immutable.
func (e *Entity) ID() uint32 {
	return uint32(*e >> 8)
}

// Next updates the Entity's Version (Generation) upon deletion in place.
func (e *Entity) next() {
	if e.Version() == 255 {
		*e -= 255
	} else {
		*e += 1
	}
}
