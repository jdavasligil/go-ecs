package queue

type Queue[T any] interface {
	Push(T)
	Pop() T
	Peek() T
	Len() int
	IsEmpty() bool
	IsFull() bool
	MemUsage() uintptr
}
