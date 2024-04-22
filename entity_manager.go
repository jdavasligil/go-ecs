package ecs

import (
	"log"

	"github.com/jdavasligil/go-ecs/pkg/queue"
)

// https://austinmorlan.com/posts/entity_component_system/#what-is-an-ecs

// The EntityManager is responsible for distributing Entity IDs.
type EntityManager struct {
	// Queue of discarded entity IDs for recycling.
	pool queue.Queue[Entity]

	// Total living entities used to enforce a limit on max entities.
	size uint32

	// ID for the next entity to be created if the recycle pool is empty.
	next Entity
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		pool:  queue.NewRingBuffer[Entity](1024),
	}
}

func (em *EntityManager) CreateEntity() Entity {
	if em.size == uint32(MAX_ENTITIES) {
		log.Printf("ECS: Failed to create entity - manager is full.")
		return 0
	}

	var id Entity

	if em.pool.IsEmpty() {
		id = em.next
		em.next += 1
	} else {
		id = em.pool.Pop()
	}

	return id
}

func (em *EntityManager) DestroyEntity(entity Entity) bool {
	if em.size == 0 {
		log.Printf("ECS: Failed to destroy entity - manager is full.")
		return false
	}

	em.pool.Push(entity)
	em.size -= 1

	return true
}
