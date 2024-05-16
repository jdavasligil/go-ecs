package ecs

import (
	"testing"

	"github.com/jdavasligil/go-ecs/pkg/testutil"
)

type myComponent struct {
	Data string
	A, B int
}

func (c myComponent) ID() ComponentID {
	return 0
}

func testComponentStore(t *testing.T) {
	store := newComponentStore[myComponent]()
	entity1 := newEntity(1)
	entity2 := newEntity(2)
	t.Run("AddRemove", func(t *testing.T) {
		store.Reset()
		testutil.AssertEqual(t, store.Size(), 0)
		testutil.AssertEqual(t, store.Add(entity1, myComponent{"A", 1, 2}), true)
		testutil.AssertEqual(t, store.Add(entity1, myComponent{"B", 2, 3}), false)
		testutil.AssertEqual(t, store.Size(), 1)
		testutil.AssertEqual(t, store.Add(entity2, myComponent{"B", 2, 3}), true)
		testutil.AssertEqual(t, store.Size(), 2)
		testutil.AssertEqual(t, store.Remove(entity1), true)
		testutil.AssertEqual(t, store.Remove(entity2), true)
		testutil.AssertEqual(t, store.Size(), 0)
	})
	t.Run("GetComponent", func(t *testing.T) {
		store.Reset()
		testutil.AssertEqual(t, store.Add(entity1, myComponent{"A", 1, 2}), true)
		c, ok := store.GetComponent(entity1)
		testutil.AssertEqual(t, ok, true)
		testutil.AssertEqual(t, c.Data, "A")
		testutil.AssertEqual(t, c.A, 1)
		testutil.AssertEqual(t, c.B, 2)
	})
	t.Run("IterateEntities", func(t *testing.T) {
		store.Reset()
		entityMap := make(map[Entity]bool)
		entities := [4]Entity{newEntity(2), newEntity(16), newEntity(1), newEntity(214)}
		values := [4]string{"A", "B", "C", "D"}
		for i, e := range entities {
			entityMap[e] = false
			testutil.AssertEqual(t, store.Add(e, myComponent{values[i], i, i + 1}), true)
		}
		store.Remove(entities[0])
		store.Add(entities[0], myComponent{"E", 1, 2})
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
		entities := [4]Entity{newEntity(2), newEntity(16), newEntity(1), newEntity(214)}
		values := [4]string{"A", "B", "C", "D"}
		for i, v := range values {
			valueMap[v] = false
			store.Add(entities[i], myComponent{v, i, i + 1})
		}
		for _, c := range store.Components() {
			valueMap[c.Data] = true
		}
		for _, v := range values {
			testutil.AssertEqual(t, valueMap[v], true)
		}
	})
}
