package ecs

import "github.com/jdavasligil/go-ecs/pkg/queue"

// https://austinmorlan.com/posts/entity_component_system/#what-is-an-ecs

type Entity uint32

type EntityManager struct {
    // Queue of unused entity IDs for recycling.
    pool queue.Queue[Entity]

    // Array of signatures where index maps to the entity ID.
    sigs []Signature

    // Total living entities used to enforce a limit on max entities.
    size uint32
}

func (em *EntityManager) NewEntityManager() *EntityManager {
    return &EntityManager{
        pool: queue.NewRingBuffer[Entity](1024),
        sigs: make([]Signature, MAX_ENTITIES),
    }
}
