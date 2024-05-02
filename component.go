package ecs

// Component represents any pure data type that is named.
type Component interface {
	TypeName() string
}
