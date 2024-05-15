package ecs_test

import "github.com/jdavasligil/go-ecs"

const (
	VelocityID ecs.ComponentID = iota
	PositionID
	HealthID
	CombatTagID
	DeadTagID
)

const (
	VelocityPageSize = 7 + iota
	PositionPageSize
	HealthPageSize
	CombatTagPageSize
	DeadTagPageSize
)

// Components
type Velocity struct {
	x float32
	y float32
	z float32
}
type Position struct {
	x float32
	y float32
	z float32
}
type Health struct {
	hp int
}

// Tags
type CombatTag struct{}
type DeadTag struct{}

// ID Methods
func (c DeadTag) ID() ecs.ComponentID {
	return DeadTagID
}
func (c Velocity) ID() ecs.ComponentID {
	return VelocityID
}
func (c Position) ID() ecs.ComponentID {
	return PositionID
}
func (c Health) ID() ecs.ComponentID {
	return HealthID
}
func (c CombatTag) ID() ecs.ComponentID {
	return CombatTagID
}
