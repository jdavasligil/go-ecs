package ecs_test

import (
	"testing"

	"github.com/jdavasligil/go-ecs"
	"github.com/jdavasligil/go-ecs/pkg/testutil"
)

type MyComponent struct {
	Data string
	A, B int
}

func (c MyComponent) ID() ecs.ComponentID {
	return 0
}

func TestComponentStore(t *testing.T) {
	store := ecs.NewComponentStore[MyComponent]()
	entity1 := ecs.NewEntity(1)
	entity2 := ecs.NewEntity(2)
	t.Run("AddRemove", func(t *testing.T) {
		store.Reset()
		testutil.AssertEqual(t, store.Size(), 0)
		testutil.AssertEqual(t, store.Add(entity1, MyComponent{"A", 1, 2}), true)
		testutil.AssertEqual(t, store.Add(entity1, MyComponent{"B", 2, 3}), false)
		testutil.AssertEqual(t, store.Size(), 1)
		testutil.AssertEqual(t, store.Add(entity2, MyComponent{"B", 2, 3}), true)
		testutil.AssertEqual(t, store.Size(), 2)
		testutil.AssertEqual(t, store.Remove(entity1), entity1)
		testutil.AssertEqual(t, store.Remove(entity2), entity2)
		testutil.AssertEqual(t, store.Size(), 0)
	})
	t.Run("GetComponent", func(t *testing.T) {
		store.Reset()
		testutil.AssertEqual(t, store.Add(entity1, MyComponent{"A", 1, 2}), true)
		c, ok := store.GetComponent(entity1)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, c.Data, "A")
		testutil.AssertEqual(t, c.A, 1)
		testutil.AssertEqual(t, c.B, 2)
	})
	t.Run("IterateEntities", func(t *testing.T) {
		store.Reset()
		entityMap := make(map[ecs.Entity]bool)
		entities := [4]ecs.Entity{ecs.NewEntity(2), ecs.NewEntity(16), ecs.NewEntity(1), ecs.NewEntity(214)}
		values := [4]string{"A", "B", "C", "D"}
		for i, e := range entities {
			entityMap[e] = false
			testutil.AssertEqual(t, store.Add(e, MyComponent{values[i], i, i + 1}), true)
		}
		store.Remove(entities[0])
		store.Add(entities[0], MyComponent{"E", 1, 2})
		for _, e := range store.Entities() {
			entityMap[e] = true
		}
		for _, e := range entities {
			testutil.AssertEqual(t, entityMap[e], true)
		}
	})
	t.Run("IterateComponents", func(t *testing.T) {
		store.Reset()
		valueMap := make(map[string]bool)
		entities := [4]ecs.Entity{ecs.NewEntity(2), ecs.NewEntity(16), ecs.NewEntity(1), ecs.NewEntity(214)}
		values := [4]string{"A", "B", "C", "D"}
		for i, v := range values {
			valueMap[v] = false
			store.Add(entities[i], MyComponent{v, i, i + 1})
		}
		for _, c := range store.Components() {
			valueMap[c.Data] = true
		}
		for _, v := range values {
			testutil.AssertEqual(t, valueMap[v], true)
		}
	})
}
