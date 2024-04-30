package pagearray

// PageArray is a space saving data structure which splits an array into pages.
// Pages initialize to nil, allowing for minimally sized gaps between pages.
type PageArray struct {
	// pages is a slice of int slices
	pages [][]int

	// Logarithmic size of page (power of 2)
	pow2 uint32
}

func NewPageArray(pow2 uint32) PageArray {
	return PageArray{
		pages: make([][]int, 0),
		pow2:  pow2,
	}
}

func (p PageArray) Set(idx int, val int) {
	pageIdx := idx >> p.pow2
	offset := idx & ((1 << p.pow2) - 1)
	for len(p.pages) <= pageIdx {
		p.pages = append(p.pages, nil)
	}
	if p.pages[pageIdx] == nil {
		p.pages[pageIdx] = make([]int, 1<<p.pow2)
		for i := 0; i < len(p.pages[len(p.pages)-1]); i++ {
			p.pages[len(p.pages)-1][i] = -1
		}
	}
	p.pages[pageIdx][offset] = val
}

func (p PageArray) Clear(idx int) {
	pageIdx := idx >> p.pow2
	offset := idx & ((1 << p.pow2) - 1)
	if len(p.pages) <= pageIdx || p.pages[pageIdx] == nil {
		return
	}
	p.pages[pageIdx][offset] = -1
}

func (p PageArray) Sweep() {
	nilOffset := 0
	for pageIdx, page := range p.pages {
		nilOffset++
		if page != nil {
			nilOffset = 0
			pageEmpty := true
			for i := 0; i < len(page); i++ {
				pageEmpty = pageEmpty && (page[i] == -1)
			}
			if pageEmpty {
				p.pages[pageIdx] = nil
			}
		}
	}
	p.pages = p.pages[:len(p.pages)-nilOffset] // TEST THIS
}

func (p PageArray) SweepAndClear(idx int) {
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

func (p PageArray) At(idx int) int {
	pageIdx := idx >> p.pow2
	offset := idx & ((1 << p.pow2) - 1)
	if len(p.pages) <= pageIdx || p.pages[pageIdx] == nil {
		return -1
	}
	return p.pages[pageIdx][offset]
}
