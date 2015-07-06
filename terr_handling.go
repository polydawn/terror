package terror

// This file contains all the guts of making terror work.
// You probably want to look at `terr.go` first for documentation and core concepts.

import (
	"bytes"
	"reflect"
	"sync/atomic"
)

/*
	Lists all the parent types in an error hierarchy.
	This shouldn't typically be necessary to look at outside this package.

	Order is from most specific to least: 0->theerror, $len-1->root.
*/
type TerrorHierarchy []reflect.Type

func (th TerrorHierarchy) String() string {
	if len(th) < 1 {
		return "-"
	}
	var buf bytes.Buffer
	for _, t := range th {
		buf.WriteString(t.Name())
		buf.WriteByte(',')
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}

func Is(this error, that Terror) bool {
	// reflect on the first member of every struct.
	// panic if it's not Terror.
	// it's the hierarchy otherwise.

	// note that one could pretty readily cache this:
	//  it'd be `map[reflect.Type][]reflect.Type`
	//   and aside from tricky first-type sync, this makes any `Is` basically the cost of the two reflects.

	// it *looks* to this reader like `reflect.TypeOf` is actually *incredibly* cheap;
	//  it's mostly just taking an unsafe pointer and wrapping it in a flat struct of accessors.
	// which seems to imply that doing equality checks on the type info should be insanely cheap.
	// i'm kind of flabbergasted by this; i expected it to be much, much costlier.
	//  if i'm reading correctly, then in fact doing hierarchies this way need be no more expensive than pointer checks the spacemonkey way.

	if this == nil {
		return false
	}
	if that == nil {
		return false
	}
	type1 := reflect.TypeOf(this)
	if type1.Kind() == reflect.Ptr {
		type1 = type1.Elem()
	}
	type2 := reflect.TypeOf(that)
	if type2.Kind() == reflect.Ptr {
		type2 = type2.Elem()
	}
	heir := GetHierarchy(type1)
	for _, t := range heir {
		if t == type2 {
			return true
		}
	}
	return false
}

/*
	Get the hierarchy information for a given type.
	This is a relatively low-level method; you probably want to use
	`Is` to compare error types instead.

	This consults the shared cache, returning the already known hierarchy if
	possible, otherwise determining it and then atomically storing it.
	(Already known hierarchies will be returned within on the order of a dozen
	nanoseconds e.g. somewhere on par with the same amount of time as a regular
	map[thing]otherthing access; initial determinations may take several
	hundred nanoseconds.  If several threads ask about a new error type
	concurrently, they may all pay the initial determination cost; this is
	generally a pretty acceptable WCS as it allows using lockless reads for
	the rest of the program, which is a significant efficiency increase over
	something with fully locked access.  Consult the benchmarks for more information.)
*/
func GetHierarchy(typ reflect.Type) TerrorHierarchy {
	hc := hierachyCacheAtom.Load().(HierarchyCache)
	th, ok := hc[typ]
	if !ok {
		//fmt.Printf("GetHierarchy hit slow path\n")
		th = DetermineHierarchy(typ)
		// make a new map (the old one is shared mem and thus should not be mutated)
		next := make(HierarchyCache)
		for k, v := range hc {
			next[k] = v // my kingdom for generics
		}
		next[typ] = th
		// put it into the atomic value; races from this same method may be lost, but
		//  since the determinations are idempotent, this is fair game.
		hierachyCacheAtom.Store(next)
	}
	return th
}

type HierarchyCache map[reflect.Type]TerrorHierarchy

var hierachyCacheAtom atomic.Value

func init() {
	hierachyCacheAtom.Store(make(HierarchyCache))
}

var terror_t reflect.Type = reflect.TypeOf((*Terror)(nil)).Elem()

func DetermineHierarchy(typ reflect.Type) TerrorHierarchy {
	hier := make(TerrorHierarchy, 0, 5)
	//fmt.Printf("determining hierarchy on type %s\n", typ)
	for i := 0; ; i++ {
		hier = append(hier, typ)
		if typ.NumField() < 1 {
			// if it doesn't have any fields, it must be a root
			//fmt.Printf("hit root: type %s doesn't have fields\n", typ)
			return hier
		}
		typ = typ.Field(0).Type
		if !typ.Implements(terror_t) {
			// if the first field isn't `Terror` type, it's a root
			//fmt.Printf("hit root: first field is type %s and doesn't implement `Terror`\n", typ)
			return hier
		}
		// if you have more than 12-deep in your error tree, you're nuts
		if i >= 12 {
			panic("error not rooted")
		}
	}
	return hier
}
