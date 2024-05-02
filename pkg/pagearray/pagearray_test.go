package pagearray_test

import (
	"testing"

	"github.com/jdavasligil/go-ecs/pkg/pagearray"
	"github.com/jdavasligil/go-ecs/pkg/testutil"
)

func TestPageArraySet(t *testing.T) {
	arr := pagearray.NewPageArray(12) // Page size: 2^12 = 4096 ints
	//t.Logf("Memory Usage (Empty): %d bytes\n", arr.MemUsage())
	arr.Set(0, 0)
	arr.Set(1, 1)
	arr.Set(4095, 4095)
	arr.Set(4096, 4096)
	cases := []struct {
		Name        string
		A, Expected int
	}{
		{"InitialValue", 0, 0},
		{"SecondValue", 1, 1},
		{"EmptyIdx2", 2, -1},
		{"BoundaryValue", 4095, 4095},
		{"BoundaryValue2", 4096, 4096},
		{"EmptyIdx4097", 4097, -1},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			testutil.AssertEqual(t, arr.At(tc.A), tc.Expected)
		})
	}
}

func TestPageArrayClear(t *testing.T) {
	arr := pagearray.NewPageArray(3) // Page size: 2^8 = 8 ints
	arr.Set(16, 0)
	testutil.AssertEqual(t, arr.At(16), 0)
	arr.Clear(16)
	testutil.AssertEqual(t, arr.At(16), -1)
}

func TestPageArraySweepAndClear(t *testing.T) {
	arr := pagearray.NewPageArray(3) // Page size: 2^8 = 8 ints
	arr.Set(1, 0)
	memInitial := arr.MemUsage()
	arr.SweepAndClear(1)
	if arr.MemUsage() >= memInitial {
		t.Errorf("MemUsage: %d, Initial: %d\n", arr.MemUsage(), memInitial)
	}
}

func TestPageArraySweep(t *testing.T) {
	arr := pagearray.NewPageArray(12) // Page size: 2^12 = 4096 ints
	memInitial := arr.MemUsage()
	arr.Set(0, 0)
	arr.Set(4096, 0)
	if arr.MemUsage() <= memInitial {
		t.Errorf("MemUsage: %d <= memInitial: %d\n", arr.MemUsage(), memInitial)
	}
	arr.Clear(0)
	arr.Clear(4096)
	arr.Sweep()
	testutil.AssertEqual(t, arr.MemUsage(), memInitial)
}

func BenchmarkPageArray(b *testing.B) {
	b.Run("PageArraySweepAndClear", func(b *testing.B) {
		pArr := pagearray.NewPageArray(12) // Page size: 2^12 = 4096 ints
		for i := 0; i < b.N; i++ {
			pArr.Set(1, 1)
			pArr.Set(4096, 2)
			pArr.Set(2, 3)
			pArr.Set(4097, 4)
			pArr.Clear(1)
			pArr.Clear(4096)
			pArr.SweepAndClear(2)
			pArr.SweepAndClear(4097)
		}
	})
	b.Run("PageArrayClear", func(b *testing.B) {
		pArr := pagearray.NewPageArray(12) // Page size: 2^12 = 4096 ints
		for i := 0; i < b.N; i++ {
			pArr.Set(1, 1)
			pArr.Set(4096, 2)
			pArr.Set(2, 3)
			pArr.Set(4097, 4)
			pArr.Clear(1)
			pArr.Clear(4096)
			pArr.Clear(2)
			pArr.Clear(4097)
		}
	})
	b.Run("Array", func(b *testing.B) {
		sArr := make([]int, 4096*2)
		for i := 0; i < b.N; i++ {
			sArr[1] = 1
			sArr[4096] = 2
			sArr[2] = 3
			sArr[4097] = 4
			sArr[1] = -1
			sArr[4096] = -1
			sArr[2] = -1
			sArr[4097] = -1
		}
	})
}
