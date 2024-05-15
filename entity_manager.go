package ecs

import (
	"log"
	"unsafe"

	"github.com/jdavasligil/go-ecs/pkg/queue"
)

// The entityManager is responsible for distributing Entity IDs.
type entityManager struct {
	MaxEntities uint32
	MaxRecycle  uint32

	// Queue of discarded entity IDs for recycling.
	bin queue.Queue[Entity]

	// Total living entities used to enforce a limit on max entities.
	size uint32

	// ID for the next entity to be created if the recycle bin is empty.
	next uint32
}

func newEntityManager(entityLimit uint32, recycleLimit uint32) entityManager {
	elim := min(MAX_ENTITIES, entityLimit)
	rlim := min(elim, recycleLimit)
	return entityManager{
		MaxEntities: elim,
		MaxRecycle:  rlim,
		bin:         queue.NewRingBuffer[Entity](int(rlim)),
		size:        0,
		next:        1,
	}
}

// Creates an entity by recycling or incrementing to the next ID.
func (em *entityManager) CreateEntity() Entity {
	if em.size == em.MaxEntities {
		log.Printf("entityManager: Failed to create entity - manager is full.")
		return 0
	}

	var entity Entity

	if em.bin.IsEmpty() {
		entity = newEntity(em.next)
		em.next += 1
	} else {
		entity = em.bin.Pop()
	}

	em.size++

	return entity
}

// RecycleEntity marks the entity as deleted and pushes it to the recycle bin.
//
// The component data must also be deleted by removing that entity from each
// associated component store handled by the component manager.
func (em *entityManager) RecycleEntity(entity Entity) bool {
	if em.size == 0 {
		log.Printf("entityManager: Failed to recycle entity - manager is empty.")
		return false
	}

	entity.Next()
	em.bin.Push(entity)
	em.size -= 1

	return true
}

func (em *entityManager) MemUsage() uintptr {
	size := unsafe.Sizeof(*em)
	size += unsafe.Sizeof(em.MaxEntities)
	size += unsafe.Sizeof(em.MaxRecycle)
	size += em.bin.MemUsage()
	size += unsafe.Sizeof(em.size)
	size += unsafe.Sizeof(em.next)
	return size
}
