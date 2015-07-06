package terror

// This file contains a hierarchicy of errors.
// These are just a prototype and suggestion; they should for now be consided pre-alpha and subject to random change.
// While there may be advantages to some basic common ground at the top of a publicly shared error hierarchy,
//  don't forget that you can freely start a new root in your own packages (all you have to do is implement `Terror`).

var _ error = E_Any{}
var _ Terror = E_Any{}

type E_Any struct{ HelpfulError }

func (E_Any) Terror()       {}
func (E_Any) Error() string { return "root error" }

type E_IO struct{ E_Any }

type E_Network struct{ E_IO }
