package terror

import (
	"reflect"
	"testing"
)

func Benchmark_TerrorGeneration(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return &E_IO{}
		case 2:
			return &E_Network{}
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		_ = err
	}
}

func Benchmark_TerrorIsaCheck(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return &E_IO{}
		case 2:
			return &E_Network{}
		default:
			panic("impossible")
		}
	}
	tmpl := &E_IO{}
	for i := 0; i < b.N; i++ {
		err := fn()
		Is(err, tmpl)
	}
}

// same as previous `Is` test, except with a deeper target (so half the answers are "no").
func Benchmark_TerrorIsaCheck2(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return &E_IO{}
		case 2:
			return &E_Network{}
		default:
			panic("impossible")
		}
	}
	tmpl := &E_Network{}
	for i := 0; i < b.N; i++ {
		err := fn()
		Is(err, tmpl)
	}
}

func Benchmark_TerrorBreakdown_Reflecting(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return &E_IO{}
		case 2:
			return &E_Network{}
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		typ := reflect.TypeOf(err)
		_ = typ
	}
}

// subtract 'Reflecting' bench to effectively find recursive reflection costs of DetermineHierarchy
func Benchmark_TerrorBreakdown_Discovering(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return E_IO{}
		case 2:
			return E_Network{}
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		if err == nil {
			continue
		}
		hier := DetermineHierarchy(reflect.TypeOf(err))
		_ = hier
	}
}

// subtract 'Reflecting' bench to effectively find cache map and synchronization costs
func Benchmark_TerrorBreakdown_Getting(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return E_IO{}
		case 2:
			return E_Network{}
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		if err == nil {
			continue
		}
		hier := GetHierarchy(reflect.TypeOf(err))
		_ = hier
	}
}

func Benchmark_TerrorBreakdown_ComparingRtypes(b *testing.B) {
	type1 := reflect.TypeOf(E_IO{})
	type2 := reflect.TypeOf(E_Network{})
	for i := 0; i < b.N; i++ {
		if type1 == type2 {
		}
	}
}

type TE0 struct{ E_IO }
type TE1 struct{ E_Network }
type TE2 struct{ E_IO }
type TE3 struct{ E_Network }
type TE4 struct{ E_IO }
type TE5 struct{ E_Network }
type TE6 struct{ E_IO }
type TE7 struct{ E_Network }
type TE8 struct{ E_IO }
type TE9 struct{ E_Network }

func Benchmark_TerrorParallelGetting(b *testing.B) {
	errs := []Terror{
		E_IO{},
		E_Network{},
		TE0{},
		TE1{},
		TE2{},
		TE3{},
		TE4{},
		TE5{},
		TE6{},
		TE7{},
		TE8{},
		TE9{},
	}
	b.SetParallelism(4)
	b.RunParallel(func(b *testing.PB) {
		var which int
		for b.Next() {
			which = (which + 1) % 12
			err := errs[which]
			hier := GetHierarchy(reflect.TypeOf(err))
			_ = hier
		}
	})
}

// TODO: use a synthetic error tree without fields, because we can really see the allocs, and would like to separate those out
// TODO: use this pattern on the last test for everyone, ffs.  it hides both allocs and func-ptr calls (which, natch, have a different and much higher cost to invoke than functions declared on a struct or package).
