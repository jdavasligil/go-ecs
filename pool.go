package ecs

// Pool is a sparse set used for each registered Component type which maps
// entities to their components.
//
// Time Complexity:
//    Add    - O(1)
//    Remove - O(n)
//    Query  - O(1)
//
type Pool[T any] struct {
    // entityIndices is a sparse array that holds the indices into EntityList.
    // The array is indexed by the entity id itself.
    // TODO: Pagination for memory conservation.
    entityIndices int

    // entityList is a packed array that contains the entities. The index
    // corresponds to the value from entityIndices.
    entityList Entity

    // componentList is a packed array that contains component data. The array
    // is aligned with entityList (i.e., entityList[i] corresponds to data in
    // componentList[i]).
    componentList T
}
