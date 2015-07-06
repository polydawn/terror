package terror

import (
	"reflect"
	"testing"
)

func TestHierarchyDetermination(t *testing.T) {
	tbl := []struct {
		Name string
		ErrT Terror
		Hier string
	}{
		{"root error", E_Any{}, "E_Any"},
		{"child error", E_IO{}, "E_IO,E_Any"},
		{"grandchild error", E_Network{}, "E_Network,E_IO,E_Any"},
	}
	for _, tt := range tbl {
		result := DetermineHierarchy(reflect.TypeOf(tt.ErrT))
		if result.String() != tt.Hier {
			t.Errorf("%q : hier should be %q, was %q", tt.Name, tt.Hier, result)
		} else {
			//t.Logf("%q : passed: hierarchy is %q", tt.Name, tt.Hier)
		}
	}
}

func TestRelationships(t *testing.T) {
	tbl := []struct {
		Name  string
		Truth bool
		One   error
		Two   Terror
	}{
		{"child isa parent", true, E_IO{}, E_Any{}},
		{"parent not isa child", false, E_IO{}, E_Network{}},
		{"grandchild isa parent", true, E_Network{}, E_Any{}},
		{"grandchild isa child", true, E_Network{}, E_IO{}},
		{"child isa itself", true, E_IO{}, E_IO{}},
		{"nil not isa parent", false, nil, E_Any{}},
		{"nil not isa nil", false, nil, nil},
		// now do them all again, with pointers for the query param; truth tables same
		{"&child isa parent", true, &E_IO{}, E_Any{}},
		{"&parent not isa child", false, &E_IO{}, E_Network{}},
		{"&grandchild isa parent", true, &E_Network{}, E_Any{}},
		{"&grandchild isa child", true, &E_Network{}, E_IO{}},
		{"&child isa itself", true, &E_IO{}, E_IO{}},
		// now again, with pointers for the guideline param; truth tables same
		{"child isa &parent", true, E_IO{}, &E_Any{}},
		{"parent not isa &child", false, E_IO{}, &E_Network{}},
		{"grandchild isa &parent", true, E_Network{}, &E_Any{}},
		{"grandchild isa &child", true, E_Network{}, &E_IO{}},
		{"child isa &itself", true, E_IO{}, &E_IO{}},
		// one more encore with pointers on both sides; truth tables same
		{"&child isa &parent", true, &E_IO{}, &E_Any{}},
		{"&parent not isa &child", false, &E_IO{}, &E_Network{}},
		{"&grandchild isa &parent", true, &E_Network{}, &E_Any{}},
		{"&grandchild isa &child", true, &E_Network{}, &E_IO{}},
		{"&child isa &itself", true, &E_IO{}, &E_IO{}},
	}
	for _, tt := range tbl {
		result := Is(tt.One, tt.Two)
		if result != tt.Truth {
			t.Errorf("%q : check should be %v, was %v", tt.Name, tt.Truth, result)
		} else {
			//t.Logf("%q : passed", tt.Name)
		}
	}
}
