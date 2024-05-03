package ecs

type ComponentID uint8

// Component represents any pure data type and is given an ID.
type Component interface {
	ID() ComponentID
}
