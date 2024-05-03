package pagearray

import (
	"unsafe"
)

// PageArray is a space saving data structure which splits an array into pages.
// Pages initialize to nil allowing for minimally sized gaps between pages.
type PageArray struct {
	pages [][]int

	// Logarithmic page size (power of 2).
	pow2 uint32
}

func NewPageArray(pow2 uint32) PageArray {
	return PageArray{
		pages: make([][]int, 0),
		pow2:  pow2,
	}
}

// Set is used to assign an index to a value. Set assumes the given index
// is non-negative. Set will grow dynamically filling the gaps with nil pages.
// The memory overhead of nil pages is small but can be cleaned up with Sweep.
func (p *PageArray) Set(idx int, val int) {
	pageIdx := idx >> p.pow2
	offset := idx & ((1 << p.pow2) - 1)
	for len(p.pages) <= pageIdx {
		p.pages = append(p.pages, nil)
	}
	if p.pages[pageIdx] == nil {
		p.pages[pageIdx] = make([]int, 1<<p.pow2)
		for i := 0; i < len(p.pages[pageIdx]); i++ {
			p.pages[pageIdx][i] = -1
		}
	}
	p.pages[pageIdx][offset] = val
}

// Clear is used to remove a value at an index by marking it as empty. Using
// Clear by itself has no effect on allocated memory. If you want to check if
// a page is empty and deallocate the unused memory, use SweepAndClear.
func (p *PageArray) Clear(idx int) {
	pageIdx := idx >> p.pow2
	offset := idx & ((1 << p.pow2) - 1)
	if len(p.pages) <= pageIdx || p.pages[pageIdx] == nil {
		return
	}
	p.pages[pageIdx][offset] = -1
}

// Sweep iterates across all pages, checks if the page is empty, and then
// deallocates that memory. It also keeps track of trailing nil pages, and
// trims them off at the end. This is a fairly expensive call and should only
// be run when necessary to clear up memory.
func (p *PageArray) Sweep() {
	nilOffset := 0
	for pageIdx, page := range p.pages {
		nilOffset++
		if page != nil {
			pageEmpty := true
			for i := 0; i < len(page); i++ {
				pageEmpty = pageEmpty && (page[i] == -1)
			}
			if pageEmpty {
				p.pages[pageIdx] = nil
			} else {
				nilOffset = 0
			}
		}
	}
	p.pages = p.pages[:len(p.pages)-nilOffset]
}

// SweepAndClear is used to remove a value at an index by marking it as empty
// and then checks if the page is empty. If the entire page is empty, then
// SweepAndClear will deallocate the whole page. The additional overhead of
// checking for empty is a quick O(N) search where N is the page size.
func (p *PageArray) SweepAndClear(idx int) {
	pageIdx := idx >> p.pow2
	offset := idx & ((1 << p.pow2) - 1)
	if len(p.pages) <= pageIdx || p.pages[pageIdx] == nil {
		return
	}
	p.pages[pageIdx][offset] = -1
	pageEmpty := true
	for i := 0; i < len(p.pages[pageIdx]); i++ {
		pageEmpty = pageEmpty && (p.pages[pageIdx][i] == -1)
	}
	if pageEmpty {
		p.pages[pageIdx] = nil
	}
}

// Reset performs a hard reset by throwing away all allocated memory for
// garbage collection. May negatively affect garbage collection performance.
func (p *PageArray) Reset() {
	p.pages = make([][]int, 0)
}

// At gets the current value at the given index, otherwise it returns -1 to
// indicate empty.
func (p *PageArray) At(idx int) int {
	pageIdx := idx >> p.pow2
	offset := idx & ((1 << p.pow2) - 1)
	// Should we do bounds checking here, or leave it up to the caller?
	if len(p.pages) <= pageIdx || p.pages[pageIdx] == nil {
		//log.Printf("LEN: %d <= IDX: %d\n", len(p.pages), pageIdx)
		return -1
	}
	return p.pages[pageIdx][offset]
}

// MemUsage returns an estimate for the current memory being used in bytes.
func (p *PageArray) MemUsage() uintptr {
	size := unsafe.Sizeof(*p) + unsafe.Sizeof(p.pages) + unsafe.Sizeof(p.pow2)
	var myInt int = 0
	for i := 0; i < len(p.pages); i++ {
		size += unsafe.Sizeof(p.pages[i])
		if p.pages[i] != nil {
			size += uintptr(cap(p.pages[i])) * unsafe.Sizeof(myInt)
		}
	}
	return size
}
