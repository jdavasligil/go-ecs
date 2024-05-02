package ecs_test

import (
	"testing"

	"github.com/jdavasligil/go-ecs"
	"github.com/jdavasligil/go-ecs/pkg/testutil"
)

type MyComponent struct {
	Name string
	A, B int
}

func TestPool(t *testing.T) {
	pool := ecs.NewPool[MyComponent]()
	entity1 := ecs.NewEntity(1)
	entity2 := ecs.NewEntity(2)
	t.Run("AddRemove", func(t *testing.T) {
		testutil.AssertEqual(t, pool.Size(), 0)
		testutil.AssertEqual(t, pool.Add(entity1, MyComponent{"A", 1, 2}), true)
		testutil.AssertEqual(t, pool.Add(entity1, MyComponent{"B", 2, 3}), false)
		testutil.AssertEqual(t, pool.Size(), 1)
		testutil.AssertEqual(t, pool.Add(entity2, MyComponent{"B", 2, 3}), true)
		testutil.AssertEqual(t, pool.Size(), 2)
		testutil.AssertEqual(t, pool.Remove(entity1), entity1)
		testutil.AssertEqual(t, pool.Remove(entity2), entity2)
		testutil.AssertEqual(t, pool.Size(), 0)
	})
	t.Run("IterateEntities", func(t *testing.T) {
		entityMap := make(map[ecs.Entity]bool)
		entities := [4]ecs.Entity{ecs.NewEntity(2), ecs.NewEntity(16), ecs.NewEntity(1), ecs.NewEntity(214)}
		values := [4]string{"A", "B", "C", "D"}
		for i, e := range entities {
			entityMap[e] = false
			testutil.AssertEqual(t, pool.Add(e, MyComponent{values[i], i, i + 1}), true)
		}
		pool.Remove(entities[0])
		pool.Add(entities[0], MyComponent{"E", 1, 2})
		for _, e := range pool.Entities() {
			entityMap[e] = true
		}
		for _, e := range entities {
			testutil.AssertEqual(t, entityMap[e], true)
		}
		pool.Reset()
	})
	t.Run("IterateComponents", func(t *testing.T) {
		valueMap := make(map[string]bool)
		entities := [4]ecs.Entity{ecs.NewEntity(2), ecs.NewEntity(16), ecs.NewEntity(1), ecs.NewEntity(214)}
		values := [4]string{"A", "B", "C", "D"}
		for i, v := range values {
			valueMap[v] = false
			pool.Add(entities[i], MyComponent{v, i, i + 1})
		}
		for _, c := range pool.Components() {
			t.Logf("Name: %s", c.Name)
			valueMap[c.Name] = true
		}
		for _, v := range values {
			testutil.AssertEqual(t, valueMap[v], true)
		}
	})
}
