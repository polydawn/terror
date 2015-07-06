package terror

import (
	"io"
	"net"
	"os"
	"testing"
)

//
// benchmarks for EXISTING STRATEGIES
//
// Benchmarks of examples of various error paradigms... using examples from
// the standard library's *various* strategies where possible.  e.g.
//   - pointer-compare errors
//   - type switching
//   - type switching with interfaces
// we won't even touch stringly typed errors.  i think everyone knows that's
// a disaster zone and shouldn't be done on pure and obvious principle.
//
// these come in pairs of "generation" functions and the real benchmark functions;
// subtract the former from the later to see the isolated cost of handling strategies
// without the costs of allocations (which, spoiler alert, predominate).
//

func Benchmark_PtrEqGeneration(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		// of course, we can't use, say `net.errTimeout`,
		//  since that's unexported :D
		// yes, this is a jibe and an admission of irritation.
		switch which {
		case 0:
			return nil
		case 1:
			return io.EOF
		case 2:
			return net.ErrWriteToConnected
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		_ = err
	}
}

func Benchmark_PtrEqSwitching(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return io.EOF
		case 2:
			return net.ErrWriteToConnected
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		switch err {
		case io.EOF:
		case net.ErrWriteToConnected:
		default:
		}
	}
}

func Benchmark_TypeGeneration(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return &os.PathError{}
		case 2:
			return &os.LinkError{}
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		_ = err
	}
}

func Benchmark_TypeSwitching(b *testing.B) {
	var which int
	// Note, first draft had these just being called in the loop...
	// but actually, allocation... well, allocation turned that from
	//  8ns/op for the PtrEq one
	//  140ns/op for this one
	// So if you're really squeezing those nanos... that *is* grounds for thought.
	//  118ns/op if you use a two-ptr empty struct
	//  24ns/op if you use a zero-field struct
	// So, anywhere under a hundred nanos (on this machine/scale), just know that you're in fret-stack-allocs territory.
	// All that said... these benches are for *handling*, not *allocation*,
	// so we're gonna ignore all that.
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return &os.PathError{}
		case 2:
			return &os.LinkError{}
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		switch err.(type) {
		case *os.PathError:
		case *os.LinkError:
		default:
		}
	}
}

func Benchmark_IfaceGeneration(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return &os.PathError{}
		case 2:
			return &net.DNSError{}
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		_ = err
	}
}

func Benchmark_IfaceSwitching(b *testing.B) {
	var which int
	fn := func() error {
		which = (which + 1) % 3
		switch which {
		case 0:
			return nil
		case 1:
			return &os.PathError{}
		case 2:
			return &net.DNSError{}
		default:
			panic("impossible")
		}
	}
	for i := 0; i < b.N; i++ {
		err := fn()
		switch err.(type) {
		case *os.PathError:
		case temporary:
		default:
		}
	}
}

// copy of what exists in stdlib `net`, which couldn't deign to export it
type temporary interface {
	Temporary() bool
}
