package terror

// This file declares basic terror structures, and should contain
//  all the docs you really need to know in order use it happily.

/*
	Terror is the marker interface for any typed error.

	Usually, if you're using Terror typed errors, you'll find it
	useful to use our hierarchy as a starting point (it'll make it easier
	to use features like stack capture).  However, if you want to use
	Terror typed errors with your own totally separate hierarchy, you can do
	that!  Just make sure your root error type implements this interface.

	Notice that this means it's possible to implement Terror semantics
	without actually linking to the `terror` package!
	As a library author, you can export error types that support Terror
	hierarchical error semantics to anyone who wants to use them that way,
	and for everyone else quietly fall back to basic type-switch behavior --
	this is frictionless for the non-Terror users, who will not be
	required to notice anything about the `terror` package or put it into
	their gopath.
*/
type Terror interface {
	Terror()
	error
}

// REVIEW: should `Terror` mix in `error`?  probably

/*
	HelpfulError is a struct for embedding that gathers all the most commonly
	helpful fields in an error -- concepts like "message", "cause" for chaining,
	stack capture, etc.

	It's optional: you can implement the `Terror` interface without it.
	But it's realllly useful.  All the built in hierarchy includes it.
	We really like stacks.

	REVIEW: I don't want to put too much weight on the Terror interface root...
	but honestly, is an error you can't even ask for a message from really that pleasant?
	I don't want to force the user to upcast constantly to do anything more complicated
	than hierarchy query.  ... though, admittedly, that *is* kinda what go does
	naturally if you're using type switches, so... maybe it's not that crazy.
*/
type HelpfulError struct {
	Message string
	Cause   *Terror // REVIEW: maybe it makes sense to have this extra type bit here, but it will require wrapper helpers if you're putting in a stdlib `error`, so it's not clear if that's gonna be super helpful.
	// TODO stacks, etc
}
