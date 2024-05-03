package ecs

import (
	"log"

	"github.com/jdavasligil/go-ecs/pkg/queue"
)

// https://austinmorlan.com/posts/entity_component_system/#what-is-an-ecs

// The EntityManager is responsible for distributing Entity IDs.
type EntityManager struct {
	// Queue of discarded entity IDs for recycling.
	bin queue.Queue[Entity]

	// Total living entities used to enforce a limit on max entities.
	size uint32

	// ID for the next entity to be created if the recycle bin is empty.
	next uint32
}

func NewEntityManager() EntityManager {
	return EntityManager{
		bin:  queue.NewRingBuffer[Entity](1024),
		size: 0,
		next: 1,
	}
}

func (em *EntityManager) CreateEntity() Entity {
	if em.size == MAX_ENTITIES {
		log.Printf("EntityManager: Failed to create entity - manager is full.")
		return 0
	}

	var entity Entity

	if em.bin.IsEmpty() {
		entity = NewEntity(em.next)
		em.next += 1
	} else {
		entity = em.bin.Pop()
	}

	return entity
}

// RecycleEntity only marks the entity as deleted and pushes it to the recycle
// bin. The component data must also be deleted by removing that entity from
// each associated component pool handled by the component manager.
func (em *EntityManager) RecycleEntity(entity Entity) bool {
	if em.size == 0 {
		log.Printf("EntityManager: Failed to recycle entity - manager is empty.")
		return false
	}

	entity.Next()
	em.bin.Push(entity)
	em.size -= 1

	return true
}
